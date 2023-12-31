package tmux

import (
	"bufio"
	"bytes"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/indeedhat/automux/internal/config"
)

// Cmd is an alias function to make running subsequent tmux commands simpler and more readable
func Cmd(conf *config.Config, parts ...string) {
	parts = append([]string{parts[0], "-t", conf.Session}, parts[1:]...)

	if conf.Debug {
		fmt.Println("tmux ", strings.Join(parts, " "))
		return
	}

	c := exec.Command("tmux", parts...)
	c.Run()
}

// SessionExists checks if there is already a tmux session with the provided session id/name
func SessionExists(conf *config.Config) bool {
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

// AwaitSession waits for the tmux session to become available before we start trying to manipulate it
func AwaitSession(c *config.Config) {
	if c.Debug {
		return
	}

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
