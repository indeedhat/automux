package cmd

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/indeedhat/automux/internal/config"
	"github.com/spf13/cobra"
)

// PrintName just prints the session name to std out
func PrintName(l *log.Logger) *cobra.Command {
	var detached bool

	cmd := &cobra.Command{
		Use:   "print-name",
		Short: "Print the session name if the target directory is a automux directory",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			configPath, err := os.Getwd()
			if err != nil {
				return err
			}

			if len(args) == 1 {
				configPath = args[0]
			}

			c, err := config.LoadAny(configPath, l, false, detached)
			if err != nil {
				return errors.New("!! invalid automux config !!\n " + err.Error())
			}

			fmt.Println(c.SessionId)

			return nil
		},
	}

	cmd.Flags().
		BoolVarP(&detached, "detached", "d", false, "Run the automux session detached\nThis will allow you to start an automux session from another session")

	return cmd
}
