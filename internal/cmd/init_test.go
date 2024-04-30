package cmd

import (
	"os"
	"testing"

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

// TestGenerateConfig makes sure that the expected template is generated for the init config
func TestGenerateConfig(t *testing.T) {
	data, err := generateConfig("tester")

	require.Nil(t, err)
	require.Equal(t, []byte(expectedGeneratedConfig), data)
}

// TestReadInput makes sure that input can be read from stdin
func TestReadInput(t *testing.T) {
	content := []byte("Hello, World!\n")
	oldStdin := os.Stdin

	tmpfile, err := os.CreateTemp("", "example")
	require.Nil(t, err, "create temp file")

	defer func() {
		tmpfile.Close()
		os.Remove(tmpfile.Name())
		os.Stdin = oldStdin
	}()

	_, err = tmpfile.Write(content)
	require.Nil(t, err, "write temp file")

	_, err = tmpfile.Seek(0, 0)
	require.Nil(t, err, "seek")

	os.Stdin = tmpfile

	data, _ := readInput()
	// require.Nil(t, err, "read err")
	require.Equal(t, content, data, "read data")
}
