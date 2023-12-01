package main

import (
	"bufio"
	"bytes"
	"errors"
	"flag"
	"fmt"
	"log"
	"strings"
	"time"

	"os"
	"os/exec"
)

const configPath = ".automux.hcl"

func main() {
	// if we are already in a tmux session then there is nothing to do
	if os.Getenv("TMUX") != "" {
		return
	}

	if _, err := os.Stat(configPath); err != nil && errors.Is(err, os.ErrNotExist) {
		return
	}

	c, err := LoadConfig(configPath)
	if err != nil {
		log.Fatal("!! invalid automux config !!\n ", err)
	}

	flag.BoolVar(&c.Debug, "debug", c.Debug, "print tmux commands rather than running them")
	flag.Parse()

	if !c.SingleSession {
		// Not totally unique as a suffix but i think good enough for this use case
		c.Session += time.Now().Format("_150405")
	}

	if sessionExists(c) {
		return
	}

	var cmd *exec.Cmd

	if c.Debug {
		fmt.Println("tmux new-session -s " + c.Session)
	} else {

		cmd = exec.Command("tmux", "new-session", "-s", c.Session)
		cmd.Stdout = os.Stdout
		cmd.Stdin = os.Stdin
		cmd.Stderr = os.Stderr
		cmd.Start()
	}

	time.Sleep(100 * time.Millisecond)
	awaitSession(c)
	processPanels(c)

	// TODO: find a way to disconnect from the session
	if !c.Debug {
		cmd.Wait()
	}
}

// processPanels walkes through the configs windows/splits an applies them to the current tmux session
func processPanels(conf *Config) {
	for i, window := range conf.Windows {
		if i != 0 {
			cmd(conf, "new-window")
		}

		// renaming the window for some reasonstops issues with blank splits
		cmd(conf, "rename-window", window.Title)

		if window.Exec != "" {
			cmd(conf, "send-keys", window.Exec, "Enter")
		}

		for _, split := range window.Splits {
			// This looks backwards but it makes the splits open in the way i expect
			orientation := "-v"
			if split.Vertical {
				orientation = "-h"
			}

			cmd(conf, "split-window", orientation)

			if split.Exec != "" {
				cmd(conf, "send-keys", split.Exec, "Enter")
			}
		}

		// stops the opening of programs from overwriting tab
		cmd(conf, "rename-window", window.Title)
	}
}

// cmd is an alias function to make running subsequent tmux commands simpler and more readable
func cmd(conf *Config, parts ...string) {
	parts = append([]string{parts[0], "-t", conf.Session}, parts[1:]...)

	if conf.Debug {
		fmt.Println("tmux ", strings.Join(parts, " "))
		return
	}

	c := exec.Command("tmux", parts...)
	c.Run()
}

// sessionExists checks if there is already a tmux session with the provided session id/name
func sessionExists(conf *Config) bool {
	c := exec.Command("tmux", "ls")
	out, err := c.CombinedOutput()
	if err != nil {
		return false
	}
	s := bufio.NewScanner(bytes.NewReader(out))

	for s.Scan() {
		parts := strings.Split(s.Text(), ":")
		if len(parts) > 0 && parts[0] == conf.Session {
			return true
		}
	}

	return false
}

// awaitSession waits for the tmux session to become available before we start trying to manipulate it
func awaitSession(c *Config) {
	ticker := time.NewTicker(50 * time.Millisecond)
	for {
		select {
		case <-time.After(time.Second):
			return
		case <-ticker.C:
			cmd := exec.Command("tmux", "ls")
			data, err := cmd.CombinedOutput()
			if err != nil {
				continue
			}

			s := bufio.NewScanner(bytes.NewReader(data))
			for s.Scan() {
				if strings.HasPrefix(s.Text(), c.Session) {
					return
				}
			}
		}
	}
}
