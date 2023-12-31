package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/indeedhat/automux/internal/config"
	"github.com/indeedhat/automux/internal/tmux"
)

// TriggerCmd is the handler for triggering the auto mux start procedure
func TriggerCmd(conf *config.Config) error {
	// if we are already in a tmux session then there is nothing to do
	if os.Getenv("TMUX") != "" {
		return nil
	}

	if !conf.SingleSession {
		// Not totally unique as a suffix but i think good enough for this use case
		conf.Session += time.Now().Format("_150405")
	}

	if tmux.SessionExists(conf) {
		return nil
	}

	cmd := exec.Command("tmux", "new-session", "-d", "-s", conf.Session)
	if conf.Config != "" {
		cmd.Args = append(cmd.Args, "-f", conf.Config)
	}

	if conf.Debug {
		fmt.Println(strings.Join(cmd.Args, " "))
	} else {
		cmd.Run()
	}

	tmux.AwaitSession(conf)
	processPanels(conf)

	if !conf.Debug {
		cmd := exec.Command("tmux", "attach", "-t", conf.Session)
		cmd.Stdout = os.Stdout
		cmd.Stdin = os.Stdin
		cmd.Stderr = os.Stderr
		cmd.Run()
		// TODO: find a way to disconnect from the session
	}
	return nil
}

// processPanels walkes through the configs windows/splits an applies them to the current tmux session
func processPanels(conf *config.Config) {
	var focus string

	for i, window := range conf.Windows {
		if window.Focus {
			focus = fmt.Sprintf(":%d.%d", i, 0)
		}

		if i != 0 {
			tmux.Cmd(conf, "new-window")
		}

		// renaming the window for some reasonstops issues with blank splits
		tmux.Cmd(conf, "rename-window", window.Title)

		if window.Exec != "" {
			tmux.Cmd(conf, "send-keys", window.Exec, "Enter")
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

			tmux.Cmd(conf, "split-window", orientation)

			if split.Size != 0 {
				tmux.Cmd(conf, "resize-pane", resize, strconv.Itoa(split.Size)+"%")
			}
			if split.Exec != "" {
				tmux.Cmd(conf, "send-keys", split.Exec, "Enter")
			}
		}

		// stops the opening of programs from overwriting tab
		tmux.Cmd(conf, "rename-window", window.Title)
	}

	if focus != "" {
		// replacing the session id is hacky and i hate it but im too lazy to come up witha proper
		// solution for now
		ses := conf.Session
		conf.Session += focus
		tmux.Cmd(conf, "select-window")
		tmux.Cmd(conf, "select-pane")
		conf.Session = ses
	}
}
