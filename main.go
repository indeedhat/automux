package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"os"
	"os/exec"
)

const configPath = ".automux.hcl"

func main() {
	var debug bool
	flag.BoolVar(&debug, "debug", false, "print tmux commands rather than running them")
	flag.Parse()

	// if we are already in a tmux session then there is nothing to do
	if os.Getenv("TMUX") != "" {
		return
	}

	if _, err := os.Stat(configPath); errors.Is(err, os.ErrNotExist) {
		return
	}

	c, err := LoadConfig(configPath)
	if err != nil {
		log.Fatal("!! invalid automux config !!\n ", err)
	}

	c.debug = debug

	if !c.SingleSession {
		// Not totally unique as a suffix but i think good enough for this use case
		c.Session += time.Now().Format("_150405")
	}

	if sessionExists(c) {
		return
	}

	cmd := exec.Command("tmux", "new-session", "-d", "-s", c.Session)
	if c.Config != "" {
		cmd.Args = append(cmd.Args, "-f", c.Config)
	}

	if c.debug {
		fmt.Println(strings.Join(cmd.Args, " "))
	} else {
		cmd.Run()
	}

	awaitSession(c)
	processPanels(c)

	if !c.debug {
		cmd := exec.Command("tmux", "attach", "-t", c.Session)
		cmd.Stdout = os.Stdout
		cmd.Stdin = os.Stdin
		cmd.Stderr = os.Stderr
		cmd.Run()
		// TODO: find a way to disconnect from the session
	}
}

// processPanels walkes through the configs windows/splits an applies them to the current tmux session
func processPanels(conf *Config) {
	var focus string

	for i, window := range conf.Windows {
		if window.Focus {
			focus = fmt.Sprintf(":%d.%d", i, 0)
		}

		if i != 0 {
			tmux(conf, "new-window")
		}

		// renaming the window for some reasonstops issues with blank splits
		tmux(conf, "rename-window", window.Title)

		if window.Exec != "" {
			tmux(conf, "send-keys", window.Exec, "Enter")
		}

		for j, split := range window.Splits {
			if split.Focus {
				focus = fmt.Sprintf(":%d.%d", i, j)
			}

			// This looks backwards but it makes the splits open in the way i expect
			orientation := "-v"
			resize := "-y"
			if split.Vertical {
				orientation = "-h"
				resize = "-x"
			}

			tmux(conf, "split-window", orientation)

			if split.Size != 0 {
				tmux(conf, "resize-pane", resize, strconv.Itoa(split.Size)+"%")
			}
			if split.Exec != "" {
				tmux(conf, "send-keys", split.Exec, "Enter")
			}
		}

		// stops the opening of programs from overwriting tab
		tmux(conf, "rename-window", window.Title)
	}

	if focus != "" {
		// replacing the session id is hacky and i hate it but im too lazy to come up witha proper
		// solution for now
		ses := conf.Session
		conf.Session += focus
		tmux(conf, "select-window")
		tmux(conf, "select-pane")
		conf.Session = ses
	}
}
