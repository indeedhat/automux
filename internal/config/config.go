package config

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/indeedhat/icl"
)

const (
	DefaultPath = ".automux"
	JsonPath    = ".automux.json"
	YamlPath    = ".automux.yml"
)

type Config struct {
	Version int `icl:"version"`
	// Used to store the relative directory for the config (if the config is not loaded from the current directory)
	Directory string
	// Session id and title for the tmux session
	SessionId string `icl:"session_id" validate:"required"`
	// AttachExisting will cause automux to re attach to any exiting session for thet directory
	AttachExisting bool `icl:"attach_existing"`
	// ConnfigPath for the tmux.conf file to use on this session
	ConfigPath string `icl:"config"`
	// Windows contains each of the tmux windo defs
	Windows []Window `icl:"window"`
	// Sessions contains definitions for background sessions to open up
	Sessions []Session `icl:"session"`

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
	Directory string `icl:".param"`

	// # Overrides:
	// Any config defined within the session block will be merged into any .automux
	// config found in the target directory with the session config taking presedence
	// over anything found there
	//
	// Session id and title for the tmux session
	SessionId string `icl:"session_id"`
	// AttachExisting will cause automux to re attach to any exiting session for thet directory
	AttachExisting *bool   `icl:"attach_existing"`
	ConfigPath     *string `icl:"config"`
	// Windows contains each of the tmux windo defs
	Windows []Window `icl:"window"`

	Debug bool
	L     *log.Logger
}

type Window struct {
	// Title of the window/tab
	Title string `icl:".param"`
	// Cmd contains the command to be run on opening the window
	Exec *string `icl:"exec"`
	// Focus sets the focus to this window after setup is done
	Focus *bool `icl:"focus"`
	// Sub directory to open the split in
	Directory *string `icl:"dir"`
	// Splits contains any extra splits to be opened in this window/tab
	Splits []Split `icl:"split"`
}

type Split struct {
	// Vertical defines if the split is vertical or horizontal
	Vertical *bool `icl:"vertical"`
	// Cmd contains any command to be ran when opening the split
	Exec *string `icl:"exec"`
	// Size in % of the total screen realestate to take up
	Size *int `icl:"size"`
	// Focus sets the focus to this split after setup is done
	Focus *bool `icl:"focus"`
	// Sub directory to open the split in
	Directory *string `icl:"dir"`
}

// Load loads the config from the given file path
func Load(path string, logger *log.Logger, debug, detached bool) (*Config, error) {
	c := Config{
		AttachExisting: true,
	}

	ast, err := icl.ParseFile(path)
	if err != nil {
		return nil, err
	} else if ast.Version() == 0 {
		return nil, errors.New(
			"you are using an old config format, see upgrade instructions" +
				"\nhttps://github.com/indeedhat/automux?tab=readme-ov-file#upgrade",
		)
	} else if ast.Version() > 1 {
		return nil, fmt.Errorf(
			"automux config version %d is not supported.\n please update automux or downgrade your config version to 1",
			ast.Version(),
		)
	}

	if err := ast.Unmarshal(&c); err != nil {
		return nil, err
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

// Exists checks if an automux config exists in the current directory
func Exists(path ...string) bool {
	p := []string{DefaultPath, JsonPath, YamlPath}
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
