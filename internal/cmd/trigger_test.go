package cmd

import (
	"bufio"
	"bytes"
	"log"
	"os"
	"strings"
	"testing"

	// "github.com/indeedhat/automux/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var triggerCmdDocument = `
version = 1
session_id = "automux-trigger-config"
window "Editor" {
	exec = "nvim"
	focus = true

	split {
		vertical = true
		exec = "htop"
		size = 20
		focus = true
	}
	split {
		size = 60
		dir = "sub/"
	}
}

session "../../_examples/" {
	session_id = "sub-automux-trigger-config-sub"
	window "Editor" {
		exec = "nvim"
		focus = true
		split {
			vertical = true
			exec = "htop"
			size = 20
			focus = true
		}
		split {
			size = 60
			dir = "sub/"
		}
	}
	window "Editor" {
		exec = "nvim"
		focus = true
		dir = "window_sub/"
		split {
			vertical = true
			exec = "htop"
			size = 20
			focus = true
		}
		split {
			size = 60
			dir = "sub/"
		}
	}
}
`

var triggerCmdDebugText = `tmux  rename-window -t automux-trigger-config Editor
tmux  send-keys -t automux-trigger-config nvim Enter
tmux  split-window -t automux-trigger-config -h
tmux  resize-pane -t automux-trigger-config -x 20%
tmux  send-keys -t automux-trigger-config htop Enter
tmux  split-window -t automux-trigger-config -v -c sub/
tmux  resize-pane -t automux-trigger-config -y 60%
tmux  rename-window -t automux-trigger-config Editor
tmux  select-window -t automux-trigger-config:0.1
tmux  select-pane -t automux-trigger-config:0.1
tmux new-session -d -s sub-automux-trigger-config-sub -c ../../_examples/
tmux  rename-window -t sub-automux-trigger-config-sub Editor
tmux  send-keys -t sub-automux-trigger-config-sub nvim Enter
tmux  split-window -t sub-automux-trigger-config-sub -h
tmux  resize-pane -t sub-automux-trigger-config-sub -x 20%
tmux  send-keys -t sub-automux-trigger-config-sub htop Enter
tmux  split-window -t sub-automux-trigger-config-sub -v -c sub/
tmux  resize-pane -t sub-automux-trigger-config-sub -y 60%
tmux  rename-window -t sub-automux-trigger-config-sub Editor
tmux  new-window -t sub-automux-trigger-config-sub -c window_sub/
tmux  rename-window -t sub-automux-trigger-config-sub Editor
tmux  send-keys -t sub-automux-trigger-config-sub nvim Enter
tmux  split-window -t sub-automux-trigger-config-sub -h -c window_sub/
tmux  resize-pane -t sub-automux-trigger-config-sub -x 20%
tmux  send-keys -t sub-automux-trigger-config-sub htop Enter
tmux  split-window -t sub-automux-trigger-config-sub -v -c window_sub/sub
tmux  resize-pane -t sub-automux-trigger-config-sub -y 60%
tmux  rename-window -t sub-automux-trigger-config-sub Editor
tmux  select-window -t sub-automux-trigger-config-sub:1.1
tmux  select-pane -t sub-automux-trigger-config-sub:1.1
`

func TestTriggerCmdTmuxSet(t *testing.T) {
	orig := os.Getenv("TMUX")
	os.Setenv("TMUX", "1")
	defer func() {
		os.Setenv("TMUX", orig)
	}()

	var b bytes.Buffer
	var l = log.New(&b, "", 0)

	assert.Nil(t, Trigger(l, ".automux").Execute(), "TriggerCmd")
}

func TestTriggerCmdMultiSession(t *testing.T) {
	os.Unsetenv("TMUX")

	tmpPath, err := os.CreateTemp("", "")
	require.Nil(t, err)
	defer os.Remove(tmpPath.Name())

	tmpPath.WriteString(triggerCmdDocument)

	var b bytes.Buffer
	var l = log.New(&b, "", 0)

	c := Trigger(l, tmpPath.Name())
	c.SetArgs([]string{"--debug", "--detached"})

	assert.Nil(t, c.Execute(), "TriggerCmd")
	parts := strings.SplitN(b.String(), "\n", 2)

	assert.True(t, strings.HasPrefix(parts[0], "tmux new-session -d -s automux-trigger-config -c /tmp/"))
	assert.Equal(t, triggerCmdDebugText, parts[1], "Debug info")
}

func t_ptr[T any](val T) *T {
	return &val
}

func t_countLines(in string) int {
	var count int

	scanner := bufio.NewScanner(strings.NewReader(in))
	for scanner.Scan() {
		count++
	}

	return count
}
