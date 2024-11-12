package config

import (
	"errors"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/hashicorp/hcl/v2/hclsimple"
)

const DefaultPath = ".automux.hcl"

type Config struct {
	// Used to store the relative directory for the config (if the config is not loaded from the current directory)
	Directory string
	// Session id and title for the tmux session
	SessionId string `hcl:"session"`
	// AttachExisting will cause automux to re attach to any exiting session for thet directory
	AttachExisting bool `hcl:"attach_existing,optional"`
	// ConnfigPath for the tmux.conf file to use on this session
	ConfigPath string `hcl:"config,optional"`
	// Windows contains each of the tmux windo defs
	Windows []Window `hcl:"window,block"`
	// Sessions contains definitions for background sessions to open up
	Sessions []Session `hcl:"session,block"`

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
	Directory string `hcl:"title,label"`

	// # Overrides:
	// Any config defined within the session block will be merged into any .automux.hcl
	// config found in the target directory with the session config taking presedence
	// over anything found there
	//
	// Session id and title for the tmux session
	SessionId string `hcl:"session,optional"`
	// AttachExisting will cause automux to re attach to any exiting session for thet directory
	AttachExisting *bool   `hcl:"attach_existing,optional"`
	ConfigPath     *string `hcl:"config,optional"`
	// Windows contains each of the tmux windo defs
	Windows []Window `hcl:"window,block"`

	Debug bool
	L     *log.Logger
}

type Window struct {
	// Title of the window/tab
	Title string `hcl:"title,label"`
	// Cmd contains the command to be run on opening the window
	Exec *string `hcl:"exec,optional"`
	// Focus sets the focus to this window after setup is done
	Focus *bool `hcl:"focus,optional"`
	// Splits contains any extra splits to be opened in this window/tab
	Splits []Split `hcl:"split,block"`
}

type Split struct {
	// Vertical defines if the split is vertical or horizontal
	Vertical *bool `hcl:"vertical,optional"`
	// Cmd contains any command to be ran when opening the split
	Exec *string `hcl:"exec,optional"`
	// Size in % of the total screen realestate to take up
	Size *int `hcl:"size,optional"`
	// Focus sets the focus to this split after setup is done
	Focus *bool `hcl:"focus,optional"`
	// Sub directory to open the split in
	Directory *string `hcl:"dir,optional"`
}

// Load loads the config from the given file path
func Load(path string, logger *log.Logger, debug, detached bool) (*Config, error) {
	c := Config{
		AttachExisting: true,
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	if err := hclsimple.Decode(path, data, nil, &c); err != nil {
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
		session.Debug = debug
		sessionConf, err := Load(filepath.Join(session.Directory, ".automux.hcl"), logger, debug, detached)
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
	p := DefaultPath
	if len(path) > 0 {
		p = path[0]
	}

	if _, err := os.Stat(p); errors.Is(err, os.ErrNotExist) {
		return false
	}

	return true
}
