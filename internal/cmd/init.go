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
	"github.com/spf13/cobra"
)

func Init() *cobra.Command {
	return &cobra.Command{
		Use:   "Init",
		Short: "Initialize automux in the current directory",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
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
		},
	}
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
