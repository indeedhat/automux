package cmd

import (
	"fmt"
	"log"

	"github.com/indeedhat/automux/internal/config"
	"github.com/spf13/cobra"
)

// PrintName just prints the session name to std out
func PrintName(l *log.Logger, configPath string) *cobra.Command {
	var detached bool

	cmd := &cobra.Command{
		Use:   "print-name",
		Short: "Print the session name if the target directory is a automux directory",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			c, err := config.Load(configPath, l, false, detached)
			if err != nil {
				log.Fatal("!! invalid automux config !!\n ", err)
			}

			fmt.Println(c.SessionId)
		},
	}

	cmd.Flags().
		BoolVarP(&detached, "detached", "d", false, "Run the automux session detached\nThis will allow you to start an automux session from another session")

	return cmd
}
