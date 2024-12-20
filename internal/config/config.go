package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/indeedhat/icl"
	"gopkg.in/yaml.v3"
)

const (
	DefaultPath = ".automux"
	JsonPath    = ".automux.json"
	YamlPath    = ".automux.yml"
	YamlAltPath = ".automux.yaml"

	defaultExt = ".automux"
	jsonExt    = ".json"
	yamlExt    = ".yml"
	yamlAltExt = ".yaml"
)

type Config struct {
	Version int `icl:"version" json:"version" yaml:"version"`
	// Used to store the relative directory for the config (if the config is not loaded from the current directory)
	Directory string
	// Session id and title for the tmux session
	SessionId string `icl:"session_id" json:"session_id" yaml:"session_id"`
	// AttachExisting will cause automux to re attach to any exiting session for thet directory
	AttachExisting bool `icl:"attach_existing" json:"attach_existing" yaml:"attach_existing"`
	// ConnfigPath for the tmux.conf file to use on this session
	ConfigPath string `icl:"config" json:"config" yaml:"config"`
	// Windows contains each of the tmux windo defs
	Windows []Window `icl:"window" json:"windows" yaml:"windows"`
	// Sessions contains definitions for background sessions to open up
	Sessions []Session `icl:"session" json:"sessions" yaml:"sessions"`

	// Cli args
	Detached bool
	Debug    bool
	L        *log.Logger
}

// AsSession converts the Config instance to a Session one
func (c *Config) AsSession() Session {
	return Session{
		Directory:      c.Directory,
		SessionId:      c.SessionId,
		AttachExisting: &c.AttachExisting,
		ConfigPath:     &c.ConfigPath,
		Windows:        c.Windows,
		Debug:          c.Debug,
		L:              c.L,
	}
}

type Session struct {
	// Directory to open the session in
	Directory string `icl:".param" json:"dir" yaml:"dir"`

	// # Overrides:
	// Any config defined within the session block will be merged into any .automux
	// config found in the target directory with the session config taking presedence
	// over anything found there
	//
	// Session id and title for the tmux session
	SessionId string `icl:"session_id" json:"session_id" yaml:"session_id"`
	// AttachExisting will cause automux to re attach to any exiting session for thet directory
	AttachExisting *bool   `icl:"attach_existing" json:"attach_existing" yaml:"attach_existing"`
	ConfigPath     *string `icl:"config" json:"config" yaml:"config"`
	// Windows contains each of the tmux windo defs
	Windows []Window `icl:"window" json:"windows" yaml:"windows"`

	Debug bool
	L     *log.Logger
}

type Window struct {
	// Title of the window/tab
	Title string `icl:".param" json:"title" yaml:"title"`
	// Cmd contains the command to be run on opening the window
	Exec *string `icl:"exec" json:"exec" yaml:"exec"`
	// Focus sets the focus to this window after setup is done
	Focus *bool `icl:"focus" json:"focus" yaml:"focus"`
	// Sub directory to open the split in
	Directory *string `icl:"dir" json:"dir" yaml:"dir"`
	// Splits contains any extra splits to be opened in this window/tab
	Splits []Split `icl:"split" json:"splits" yaml:"splits"`
}

type Split struct {
	// Vertical defines if the split is vertical or horizontal
	Vertical *bool `icl:"vertical" json:"vertical" yaml:"vertical"`
	// Cmd contains any command to be ran when opening the split
	Exec *string `icl:"exec" json:"exec" yaml:"exec"`
	// Size in % of the total screen realestate to take up
	Size *int `icl:"size" json:"size" yaml:"size"`
	// Focus sets the focus to this split after setup is done
	Focus *bool `icl:"focus" json:"focus" yaml:"focus"`
	// Sub directory to open the split in
	Directory *string `icl:"dir" json:"dir" yaml:"dir"`
}

// Exists checks if an automux config exists in the current directory
func Exists(path ...string) bool {
	p := []string{DefaultPath, JsonPath, YamlPath, YamlAltPath}
	if len(path) > 0 {
		p = path
	}

	for _, path := range p {
		if _, err := os.Stat(path); err == nil {
			return true
		}
	}

	return false
}

// LoadAny loads the first available config from the provided dir
func LoadAny(path string, logger *log.Logger, debug, detached bool) (*Config, error) {
	stat, err := os.Stat(path)
	if err != nil || !stat.IsDir() {
		return Load(path, logger, debug, detached)
	}

	path = filepath.Join(path, defaultExt)

	if c, err := Load(path, logger, debug, detached); err == nil {
		return c, nil
	}

	if c, err := Load(path+jsonExt, logger, debug, detached); err == nil {
		return c, nil
	}

	if c, err := Load(path+yamlExt, logger, debug, detached); err == nil {
		return c, nil
	}

	if c, err := Load(path+yamlAltExt, logger, debug, detached); err == nil {
		return c, nil
	}

	return nil, os.ErrNotExist
}

// Load loads the config from the given file path
func Load(path string, logger *log.Logger, debug, detached bool) (*Config, error) {
	c := Config{
		AttachExisting: true,
	}

	ext := filepath.Ext(path)

	switch ext {
	case defaultExt:
		if err := loadICL(path, &c); err != nil {
			return nil, err
		}
	case jsonExt:
		if err := loadJSON(path, &c); err != nil {
			return nil, err
		}
	case yamlExt, yamlAltExt:
		if err := loadYAML(path, &c); err != nil {
			return nil, err
		}
	default:
		return nil, errors.New("Config not found")
	}

	// stop spaces from breaking the tmux commands
	c.SessionId = strings.ReplaceAll(c.SessionId, " ", "-")
	c.Debug = debug
	c.Detached = detached
	c.L = logger

	if path != DefaultPath {
		c.Directory = strings.TrimSuffix(path, DefaultPath)
	}

	var validSessions []Session
	for _, session := range c.Sessions {
		session.L = logger
		session.Debug = debug
		sessionConf, err := Load(filepath.Join(session.Directory, ".automux"), logger, debug, detached)
		if err != nil {
			if os.IsNotExist(err) {
				validSessions = append(validSessions, session)
			}
			continue
		}

		validSessions = append(validSessions, mergeSessions(sessionConf.AsSession(), session))
	}

	c.Sessions = validSessions

	return &c, nil
}

func loadICL(path string, c *Config) error {
	ast, err := icl.ParseFile(path)
	if err != nil {
		return err
	}

	if err := versionCheck(ast.Version()); err != nil {
		return err
	}

	return ast.Unmarshal(c)
}

func loadJSON(path string, c *Config) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(data, c); err != nil {
		return err
	}

	return versionCheck(c.Version)
}

func loadYAML(path string, c *Config) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	if err := yaml.Unmarshal(data, c); err != nil {
		return err
	}

	return versionCheck(c.Version)
}

func versionCheck(version int) error {
	if version == 0 {
		return errors.New(
			"you are using an old config format, see upgrade instructions" +
				"\nhttps://github.com/indeedhat/automux?tab=readme-ov-file#upgrade",
		)
	} else if version > 1 {
		return fmt.Errorf(
			"automux config version %d is not supported.\n please update automux or downgrade your config version to 1",
			version,
		)
	}

	return nil
}
