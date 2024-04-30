package tmux

import (
	"os/exec"
	"testing"
	"time"

	"github.com/indeedhat/automux/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var sessionExistsChecks = []struct {
	id       string
	expected bool
}{
	{"automux-test-session", true},
	{"automux-not-exists", false},
}

// TestSessionExists checks both paths for session existance
func TestSessionExists(t *testing.T) {
	s := config.Session{SessionId: "automux-test-session"}

	c := exec.Command("tmux", "new-session", "-d", "-s", s.SessionId)
	require.Nil(t, c.Run(), "setup session")

	for _, check := range sessionExistsChecks {
		t.Run(check.id, func(t *testing.T) {
			s.SessionId = check.id
			require.Equal(t, check.expected, SessionExists(s), "session existis")
		})
	}

	c = exec.Command("tmux", "kill-session", "-t", "automux-test-session")
	require.Nil(t, c.Run(), "kill session")
}

// TestAwaitSession checks that after a session is started AwaitSession will find the session before it hits timeout
func TestAwaitSession(t *testing.T) {
	s := config.Session{SessionId: "automux-test-session"}

	// setup session to test against
	c := exec.Command("tmux", "new-session", "-d", "-s", s.SessionId)
	require.Nil(t, c.Run(), "setup session")

	start := time.Now()
	AwaitSession(s)
	assert.WithinDuration(t, start, time.Now(), time.Second, "found session")

	c = exec.Command("tmux", "kill-session", "-t", s.SessionId)
	require.Nil(t, c.Run(), "kill session")
}

// TestAwaitSessionTimeout checks that if no session is found AwaitSession will timeout after a second
func TestAwaitSessionTimeout(t *testing.T) {
	s := config.Session{SessionId: "automux-test-session"}

	start := time.Now()
	AwaitSession(s)

	assert.WithinDuration(t, start, time.Now(), time.Millisecond*1015, "times out soon after a seccond")
}

// TestAwaitSessionDebug checks that AwaitSession does not actully wait for a session in debug mode
func TestAwaitSessionDebug(t *testing.T) {
	s := config.Session{SessionId: "automux-test-session", Debug: true}

	start := time.Now()
	AwaitSession(s)

	assert.WithinDuration(t, start, time.Now(), time.Millisecond*10, "times out soon after a seccond")
}
