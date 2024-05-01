package cmd

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var expectedGeneratedConfig = `session = "tester"
# config = "./tmux.conf"
# single_session = false

window "Editor" {
    exec = "vim"
    focus = true

    # split {
    #     vertical = true
    #     exec = "cmd_to_run_in_split"
    #     size = 30
    #     vertical = true
    # }
}

window "Shell" {
    split {}
}
`

func TestInitCmd(t *testing.T) {
	tmpStdin, stdin := t_setupStdin(t)
	defer func() {
		tmpStdin.Close()
		os.Remove(tmpStdin.Name())
		os.Stdin = stdin
		os.Remove(".automux.hcl")
	}()

	assert.Nil(t, InitCmd(), "initCmd")
	assert.FileExists(t, ".automux.hcl")

	stat, err := os.Stat(".automux.hcl")
	assert.Nil(t, err, "stat")
	modTime := stat.ModTime()

	assert.Nil(t, InitCmd(), "initCmd")
	stat, err = os.Stat(".automux.hcl")
	assert.Nil(t, err, "stat")
	assert.Equal(t, modTime, stat.ModTime(), "file not updated")
}

func t_setupStdin(t *testing.T) (*os.File, *os.File) {
	content := []byte("Hello, World!\n")
	oldStdin := os.Stdin

	tmpfile, err := os.CreateTemp("", "example")
	require.Nil(t, err, "create temp file")

	_, err = tmpfile.Write(content)
	require.Nil(t, err, "write temp file")

	_, err = tmpfile.Seek(0, 0)
	require.Nil(t, err, "seek")

	os.Stdin = tmpfile

	return tmpfile, oldStdin
}
