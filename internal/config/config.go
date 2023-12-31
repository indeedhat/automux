package config

import (
	"errors"
	"os"
	"strings"

	"github.com/hashicorp/hcl/v2/hclsimple"
)

const DefaultPath = ".automux.hcl"

type Config struct {
	// Session id and title for the tmux session
	Session string `hcl:"session"`
	// SingleSession when set automux will not run if there is already a tmux session with the
	// provided {session}
	SingleSession bool   `hcl:"single_session,optional"`
	Config        string `hcl:"config,optional"`
	// Windows contains each of the tmux windo defs
	Windows []Window `hcl:"window,block"`

	// Cli args
	Debug bool
}

type Window struct {
	// Title of the window/tab
	Title string `hcl:"title,label"`
	// Cmd contains the command to be run on opening the window
	Exec string `hcl:"exec,optional"`
	// Splits contains any extra splits to be opened in this window/tab
	Splits []Split `hcl:"split,block"`
	// Focus sets the focus to this window after setup is done
	Focus bool `hcl:"focus,optional"`
}

type Split struct {
	// Vertical defines if the split is vertical or horizontal
	Vertical bool `hcl:"vertical,optional"`
	// Cmd contains any command to be ran when opening the split
	Exec string `hcl:"exec,optional"`
	// Size in % of the total screen realestate to take up
	Size int `hcl:"size,optional"`
	// Focus sets the focus to this split after setup is done
	Focus bool `hcl:"focus,optional"`
}

// Load loads the config from the given file path
func Load(path string) (*Config, error) {
	var c Config

	if err := hclsimple.DecodeFile(path, nil, &c); err != nil {
		return nil, err
	}

	// stop spaces from breaking the tmux commands
	c.Session = strings.ReplaceAll(c.Session, " ", "-")

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
