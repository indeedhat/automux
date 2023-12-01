package main

import (
	"strings"

	"github.com/hashicorp/hcl/v2/hclsimple"
)

type Config struct {
	// Session id and title for the tmux session
	Session string `hcl:"session"`
	// SingleSession when set automux will not run if there is already a tmux session with the
	// provided {session}
	SingleSession bool `hcl:"single_session,optional"`
	// Windows contains each of the tmux windo defs
	Windows []Window `hcl:"window,block"`

	// Cli args
	debug bool
}

type Window struct {
	// Title of the window/tab
	Title string `hcl:"title,label"`
	// Cmd contains the command to be run on opening the window
	Exec string `hcl:"exec,optional"`
	// Splits contains any extra splits to be opened in this window/tab
	Splits []Split `hcl:"split,block"`
}

type Split struct {
	// Vertical defines if the split is vertical or horizontal
	Vertical bool `hcl:"vertical,optional"`
	// Cmd contains any command to be ran when opening the split
	Exec string `hcl:"exec,optional"`
	// Size in % of the total screen realestate to take up
	Size int `hcl:"size,optional"`
}

// LoadConfig loads the config from the given file path
func LoadConfig(path string) (*Config, error) {
	var c Config

	if err := hclsimple.DecodeFile(path, nil, &c); err != nil {
		return nil, err
	}

	// stop spaces from breaking the tmux commands
	c.Session = strings.ReplaceAll(c.Session, " ", "-")

	return &c, nil
}
