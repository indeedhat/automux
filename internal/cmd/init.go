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
	var jsonFlag bool
	var yamlFlag bool

	cmd := &cobra.Command{
		Use:   "init",
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

			tpl := configs.IclTemplate
			path := config.DefaultPath
			if jsonFlag {
				tpl = configs.JsonTemplate
				path = config.JsonPath
			} else if yamlFlag {
				tpl = configs.YamlTemplate
				path = config.YamlPath
			}

			configTpl, err := generateConfig(tpl, string(input))
			if err != nil {
				return err
			}

			if err = os.WriteFile(path, configTpl, 0644); err == nil {
				fmt.Print("AutoMux config created\n")
			}

			return err
		},
	}

	cmd.Flags().BoolVar(&jsonFlag, "json", false, "Create the config in json format")
	cmd.Flags().BoolVar(&yamlFlag, "yaml", false, "Create the config in yaml format")

	return cmd
}

// readInput reads a single line of input from stdin
func readInput() ([]byte, error) {
	reader := bufio.NewReader(os.Stdin)
	return reader.ReadBytes('\n')
}

// generateConfig writes the default config file to the current directory
func generateConfig(t, name string) ([]byte, error) {
	tmpl, err := template.New("config").Parse(t)
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
