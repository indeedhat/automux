package cmd

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strings"
	"text/template"

	"github.com/indeedhat/automux/configs"
	"github.com/indeedhat/automux/internal/config"
)

// InitCmd handles setting up a default automux config in the current directory
func InitCmd() error {
	if config.Exists() {
		return nil
	}

	fmt.Print("Enter the session name: ")
	input, err := readInput()
	if err != nil {
		return err
	}

	configTpl, err := generateConfig(string(input))
	if err != nil {
		return err
	}

	if err = os.WriteFile(config.DefaultPath, configTpl, 0644); err == nil {
		fmt.Print("AutoMux config created\n")
	}

	return err
}

// readInput reads a single line of input from stdin
func readInput() ([]byte, error) {
	reader := bufio.NewReader(os.Stdin)
	return reader.ReadBytes('\n')
}

// generateConfig writes the default config file to the current directory
func generateConfig(name string) ([]byte, error) {
	tmpl, err := template.New("config").Parse(configs.ConfigTemplate)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer

	name = strings.ReplaceAll(name, " ", "-")
	name = strings.ReplaceAll(name, "\r", "")
	name = strings.ReplaceAll(name, "\n", "")

	if err = tmpl.Execute(&buf, struct{ SessionName string }{name}); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
