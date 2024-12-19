package cmd

import (
	"os"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/require"
)

var expectedIclConfig = `# github.com/indeedhat/automux
# config version
version = 1

session_id = "tester"
# config = "./tmux.conf"
# attach_existing = false

window "Editor" {
    exec = "vim"
    focus = true
}

window "Shell" {
    split {
        #     vertical = true
        #     exec = "cmd_to_run_in_split"
        #     size = 30
        #     dir = "sub/"
    }
}

# vi: ft=hcl
`

func TestInitCmd(t *testing.T) {
	tmpStdin, stdin := t_setupStdin(t)
	defer func() {
		tmpStdin.Close()
		os.Remove(tmpStdin.Name())
		os.Stdin = stdin
		os.Remove(".automux")
	}()

	require.Nil(t, InitC().Execute(), "initCmd")
	require.FileExists(t, ".automux")

	stat, err := os.Stat(".automux")
	require.Nil(t, err, "stat")
	modTime := stat.ModTime()

	require.Nil(t, InitC().Execute(), "initCmd")
	stat, err = os.Stat(".automux")
	require.Nil(t, err, "stat")
	require.Equal(t, modTime, stat.ModTime(), "file not updated")

	data, err := os.ReadFile(".automux")
	require.Nil(t, err)
	require.Equal(t, expectedIclConfig, string(data))
}

var expectedJsonConfig = `{
  "version": 1,
  "session_id": "tester",
  "windows": [
    {
      ".title": "Editor",
      "exec": "vim",
      "focus": true
    },
    {
      ".title": "Shell",
      "splits": [
        {}
      ]
    }
  ]
}
`

func TestInitCmdWithJson(t *testing.T) {
	tmpStdin, stdin := t_setupStdin(t)
	defer func() {
		tmpStdin.Close()
		os.Remove(tmpStdin.Name())
		os.Stdin = stdin
		os.Remove(".automux.json")
	}()

	c := InitC()
	c.SetArgs([]string{"--json"})

	require.Nil(t, c.Execute(), "initCmd")
	require.FileExists(t, ".automux.json")

	stat, err := os.Stat(".automux.json")
	require.Nil(t, err, "stat")
	modTime := stat.ModTime()

	c = InitC()
	c.SetArgs([]string{"--json"})
	spew.Dump(c)
	require.Nil(t, c.Execute(), "initCmd")

	stat, err = os.Stat(".automux.json")
	require.Nil(t, err, "stat")
	require.Equal(t, modTime, stat.ModTime(), "file not updated")

	data, err := os.ReadFile(".automux.json")
	require.Nil(t, err)
	require.Equal(t, expectedJsonConfig, string(data))
}

var expectedYamlConfig = `# config version
version: 1
session_id: "tester"
windows:
  - .title: Editor
    exec: vim
    focus: true
  - .title: Shell
    splits:
      - {}
`

func TestInitCmdWithYaml(t *testing.T) {
	tmpStdin, stdin := t_setupStdin(t)
	defer func() {
		tmpStdin.Close()
		os.Remove(tmpStdin.Name())
		os.Stdin = stdin
		os.Remove(".automux.yml")
	}()

	c := InitC()
	c.SetArgs([]string{"--yaml"})

	require.NoFileExists(t, ".automux")
	require.Nil(t, c.Execute(), "initCmd")
	require.FileExists(t, ".automux.yml")

	stat, err := os.Stat(".automux.yml")
	require.Nil(t, err, "stat")
	modTime := stat.ModTime()

	c = InitC()
	c.SetArgs([]string{"--yaml"})
	require.Nil(t, c.Execute(), "initCmd")

	stat, err = os.Stat(".automux.yml")
	require.Nil(t, err, "stat")
	require.Equal(t, modTime, stat.ModTime(), "file not updated")

	data, err := os.ReadFile(".automux.yml")
	require.Nil(t, err)
	require.Equal(t, expectedYamlConfig, string(data))
}

func t_setupStdin(t *testing.T) (*os.File, *os.File) {
	content := []byte("tester\n")
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
