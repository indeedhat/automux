package cmd

import (
	"bufio"
	"bytes"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/indeedhat/automux/internal/config"
	"github.com/stretchr/testify/assert"
)

var triggerCmdConfig = &config.Config{
	SessionId:  "automux-trigger-config",
	ConfigPath: ".automux.hcl",
	Windows: []config.Window{
		{
			Title: "Editor",
			Exec:  t_ptr("nvim"),
			Splits: []config.Split{
				{
					Vertical: t_ptr(true),
					Exec:     t_ptr("htop"),
					Size:     t_ptr(20),
					Focus:    t_ptr(true),
				},
				{
					Size: t_ptr(60),
				},
			},
		},
		{
			Title: "Extras",
		},
	},
	Debug: true,
	Sessions: []config.Session{
		{
			Debug:     true,
			Directory: "../../_examples/single_session",
			SessionId: "sub-automux-trigger-config-sub",
			Windows: []config.Window{
				{
					Title: "Editor",
					Exec:  t_ptr("nvim"),
					Focus: t_ptr(true),
					Splits: []config.Split{
						{
							Vertical: t_ptr(true),
							Exec:     t_ptr("htop"),
							Size:     t_ptr(20),
							Focus:    t_ptr(true),
						},
						{
							Size:      t_ptr(60),
							Directory: t_ptr("sub/"),
						},
					},
				},
				{
					Title:     "Editor",
					Exec:      t_ptr("nvim"),
					Focus:     t_ptr(true),
					Directory: t_ptr("window_sub/"),
					Splits: []config.Split{
						{
							Vertical: t_ptr(true),
							Exec:     t_ptr("htop"),
							Size:     t_ptr(20),
							Focus:    t_ptr(true),
						},
						{
							Size:      t_ptr(60),
							Directory: t_ptr("sub/"),
						},
					},
				},
			},
		},
	},
}

var triggerCmdDebugText = `tmux new-session -d -s automux-trigger-config -f .automux.hcl
tmux  rename-window -t automux-trigger-config Editor
tmux  send-keys -t automux-trigger-config nvim Enter
tmux  split-window -t automux-trigger-config -h
tmux  resize-pane -t automux-trigger-config -x 20%
tmux  send-keys -t automux-trigger-config htop Enter
tmux  split-window -t automux-trigger-config -v
tmux  resize-pane -t automux-trigger-config -y 60%
tmux  rename-window -t automux-trigger-config Editor
tmux  new-window -t automux-trigger-config
tmux  rename-window -t automux-trigger-config Extras
tmux  rename-window -t automux-trigger-config Extras
tmux  select-window -t automux-trigger-config:0.0
tmux  select-pane -t automux-trigger-config:0.0
tmux new-session -d -s sub-automux-trigger-config-sub -c ../../_examples/single_session
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
tmux  select-window -t sub-automux-trigger-config-sub:1.0
tmux  select-pane -t sub-automux-trigger-config-sub:1.0
`

func TestTriggerCmdTmuxSet(t *testing.T) {
	orig := os.Getenv("TMUX")
	os.Setenv("TMUX", "1")
	defer func() {
		os.Setenv("TMUX", orig)
	}()

	assert.Nil(t, TriggerCmd(triggerCmdConfig), "TriggerCmd")
}

func TestTriggerCmdDebug(t *testing.T) {
	os.Unsetenv("TMUX")

	var b bytes.Buffer
	triggerCmdConfig.L = log.New(&b, "", 0)
	triggerCmdConfig.Sessions[0].L = triggerCmdConfig.L

	assert.Nil(t, TriggerCmd(triggerCmdConfig), "TriggerCmd")
	assert.Equal(t, triggerCmdDebugText, b.String(), "Debug info")
}

func TestTriggerCmdMultiSession(t *testing.T) {
	os.Unsetenv("TMUX")

	var b bytes.Buffer
	triggerCmdConfig.L = log.New(&b, "", 0)
	triggerCmdConfig.Sessions[0].L = triggerCmdConfig.L

	assert.Nil(t, TriggerCmd(triggerCmdConfig), "TriggerCmd")
	// this is simpler than rebuilding the output data from a template with the current session suffix
	assert.Equal(t, t_countLines(triggerCmdDebugText), t_countLines(b.String()), "Debug info")
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
