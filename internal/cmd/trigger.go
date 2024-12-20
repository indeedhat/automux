package cmd

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"

	"github.com/indeedhat/automux/internal/config"
	"github.com/indeedhat/automux/internal/tmux"
	"github.com/spf13/cobra"
)

var (
	triggerFlagDebug    bool
	triggerFlagDetached bool
)

func Trigger() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "",
		Short: "Trigger the automux config in the current directory, if present",
		Args:  cobra.MaximumNArgs(1),
		RunE:  triggerCmd,
	}

	cmd.Flags().BoolVar(&triggerFlagDebug, "debug", false, "print tmux commands rather than running them")
	cmd.Flags().BoolVarP(
		&triggerFlagDetached,
		"detached",
		"d",
		false,
		"Run the automux session detached\nThis will allow you to start an automux session from another session",
	)

	return cmd
}

func triggerCmd(cmd *cobra.Command, args []string) error {
	// if we are already in a tmux session then there is nothing to do
	if os.Getenv("TMUX") != "" && !triggerFlagDetached {
		return nil
	}

	configPath, err := os.Getwd()
	if err != nil {
		return err
	}

	if len(args) == 1 {
		configPath = args[0]
	}

	conf, err := config.LoadAny(
		configPath,
		cmd.Context().Value("logger").(*log.Logger),
		triggerFlagDebug,
		triggerFlagDetached,
	)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}

		return errors.New("!! invalid automux config !!\n " + err.Error())
	}

	masterSession := conf.AsSession()

	if tmux.SessionExists(masterSession) {
		if conf.AttachExisting {
			goto attach
		}

		return nil
	}

	createSession(masterSession)

	for i, session := range conf.Sessions {
		if session.SessionId == "" {
			conf.L.Printf("Failed to start session %d: no session id set\n", i)
			continue
		}
		if tmux.SessionExists(session) {
			continue
		}

		createSession(session)
	}

attach:
	if !conf.Debug && !conf.Detached {
		cmd := exec.Command("tmux", "attach", "-t", conf.SessionId)
		cmd.Stdout = os.Stdout
		cmd.Stdin = os.Stdin
		cmd.Stderr = os.Stderr
		cmd.Run()
	}

	return nil
}

// createSession creates a new tmux session, wait for the server to start it then
// create the sessions layout based on the provided config
func createSession(session config.Session) {
	args := []string{"new-session", "-d", "-s", session.SessionId}
	if session.Directory != "" {
		args = append(args, "-c", session.Directory)
	}

	cmd := exec.Command("tmux", args...)
	if session.ConfigPath != nil && *session.ConfigPath != "" {
		cmd.Args = append(cmd.Args, "-f", *session.ConfigPath)
	}

	if session.Directory != "" {
		cmd.Dir = session.Directory
	}

	if session.Debug {
		session.L.Println(strings.Join(cmd.Args, " "))
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
			if window.Directory != nil && *window.Directory != "" {
				tmux.Cmd(session, "new-window", "-c", *window.Directory)
			} else {
				tmux.Cmd(session, "new-window")
			}
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
			*focus = fmt.Sprintf(":%d.%d", i, j+1)
		}

		// This looks backwards but it makes the splits open in the way i expect
		orientation := "-v"
		resize := "-y"
		if split.Vertical != nil && *split.Vertical {
			orientation = "-h"
			resize = "-x"
		}

		if window.Directory != nil && *window.Directory != "" {
			if split.Directory != nil {
				*split.Directory = path.Join(*window.Directory, *split.Directory)
			} else {
				split.Directory = window.Directory
			}
		}

		splitArgs := []string{"split-window", orientation}
		if split.Directory != nil && *split.Directory != "" {
			splitArgs = append(splitArgs, "-c", *split.Directory)
		}

		tmux.Cmd(session, splitArgs...)

		if split.Size != nil && *split.Size != 0 {
			tmux.Cmd(session, "resize-pane", resize, strconv.Itoa(*split.Size)+"%")
		}
		if split.Exec != nil && *split.Exec != "" {
			tmux.Cmd(session, "send-keys", *split.Exec, "Enter")
		}
	}
}
