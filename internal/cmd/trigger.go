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
		conf.SessionId += time.Now().Format("_150405")
	}

	masterSession := conf.AsSession()

	if tmux.SessionExists(masterSession) {
		return nil
	}

	createSession(masterSession)

	for i, session := range conf.Sessions {
		if session.SessionId == "" {
			fmt.Printf("Failed to start session %d: no session id set\n", i)
			continue
		}
		if session.SingleSession != nil && !*session.SingleSession {
			session.SessionId += time.Now().Format("_150405")
		}
		if tmux.SessionExists(session) {
			continue
		}

		createSession(session)
	}

	if !conf.Debug {
		cmd := exec.Command("tmux", "attach", "-t", conf.SessionId)
		cmd.Stdout = os.Stdout
		cmd.Stdin = os.Stdin
		cmd.Stderr = os.Stderr
		cmd.Run()
		// TODO: find a way to disconnect from the session
	}
	return nil
}

// createSession creates a new tmux session, wait for the server to start it then
// create the sessions layout based on the provided config
func createSession(session config.Session) {
	cmd := exec.Command("tmux", "new-session", "-d", "-s", session.SessionId)
	if session.ConfigPath != nil && *session.ConfigPath != "" {
		cmd.Args = append(cmd.Args, "-f", *session.ConfigPath)
	}

	if session.Directory != "" {
		cmd.Dir = session.Directory
	}

	if session.Debug {
		fmt.Println(strings.Join(cmd.Args, " "))
	} else {
		cmd.Run()
	}

	tmux.AwaitSession(session)
	processPanels(session)
}

// processPanels walkes through the configs windows/splits an applies them to the current tmux session
func processPanels(session config.Session) {
	var focus string

	for i, window := range session.Windows {
		if window.Focus != nil && *window.Focus {
			focus = fmt.Sprintf(":%d.%d", i, 0)
		}

		if i != 0 {
			tmux.Cmd(session, "new-window")
		}

		// renaming the window for some reasonstops issues with blank splits
		tmux.Cmd(session, "rename-window", window.Title)

		if window.Exec != nil && *window.Exec != "" {
			tmux.Cmd(session, "send-keys", *window.Exec, "Enter")
		}

		processSplits(window, session, &focus, i)

		// stops the opening of programs from overwriting tab
		tmux.Cmd(session, "rename-window", window.Title)
	}

	if focus != "" {
		// replacing the session id is hacky and i hate it but im too lazy to come up witha proper
		// solution for now
		ses := session.SessionId
		session.SessionId += focus
		tmux.Cmd(session, "select-window")
		tmux.Cmd(session, "select-pane")
		session.SessionId = ses
	}
}

// processSplits loops over the windows splits and adds them to the session
func processSplits(window config.Window, session config.Session, focus *string, i int) {
	for j, split := range window.Splits {
		if split.Focus != nil && *split.Focus {
			*focus = fmt.Sprintf(":%d.%d", i, j)
		}

		// This looks backwards but it makes the splits open in the way i expect
		orientation := "-v"
		resize := "-y"
		if split.Vertical != nil && *split.Vertical {
			orientation = "-h"
			resize = "-x"
		}

		tmux.Cmd(session, "split-window", orientation)

		if split.Size != nil && *split.Size != 0 {
			tmux.Cmd(session, "resize-pane", resize, strconv.Itoa(*split.Size)+"%")
		}
		if split.Exec != nil && *split.Exec != "" {
			tmux.Cmd(session, "send-keys", *split.Exec, "Enter")
		}
	}
}
