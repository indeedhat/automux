package config

import (
	"bytes"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestConfigAsSession ensures that the correct values are carried over when calling AsSession on &Config
func TestConfigAsSession(t *testing.T) {
	c := &Config{
		SessionId:     "automux-test-session",
		SingleSession: true,
		ConfigPath:    "some/fake/path/.automux",
		Windows: []Window{
			{Title: "automux-test-title"},
		},
		Debug: true,
	}

	s := c.AsSession()

	require.Equal(t, c.SessionId, s.SessionId)
	require.Equal(t, c.SingleSession, *s.SingleSession)
	require.Equal(t, c.ConfigPath, *s.ConfigPath)
	require.Equal(t, c.Windows[0].Title, s.Windows[0].Title)
	require.Equal(t, c.Debug, s.Debug)
}

var loadChecks = []struct {
	name          string
	path          string
	debug         bool
	shouldSucceed bool
	sessionCount  int
}{
	{"single-session", "../../_examples/single_session", true, true, 0},
	{"multi-session", "../../_examples/multi_session", true, true, 3},
	{"bad-path", "../../_examples", false, false, 0},
	{"bad-config", "../../_examples/bad-config", false, false, 0},
}

// TestLoad checks that config is correctly loaded by file path
func TestLoad(t *testing.T) {
	testDir, err := os.Getwd()
	require.Nil(t, err)
	var (
		b bytes.Buffer
		l = log.New(&b, "", 0)
	)

	for _, check := range loadChecks {
		t.Run(check.name, func(t *testing.T) {
			require.Nil(t, os.Chdir(check.path))

			c, err := Load(".automux.hcl", l, check.debug, false)
			if !check.shouldSucceed {
				require.NotNil(t, err)
				return
			}

			require.Nil(t, err)
			require.NotNil(t, c)
			require.Len(t, c.Sessions, check.sessionCount)
		})

		require.Nil(t, os.Chdir(testDir))
	}
}

var existsChecks = []struct {
	name     string
	path     []string
	expected bool
}{
	{"current-dir", []string{}, false},
	{"exists", []string{"../../_examples/single_session"}, true},
	{"not-exists", []string{"../../_nope/exists"}, false},
}

// TestExists makes sure that the Exists function actually works
func TestExists(t *testing.T) {
	for _, check := range existsChecks {
		t.Run(check.name, func(t *testing.T) {
			require.Equal(t, check.expected, Exists(check.path...))
		})
	}
}
