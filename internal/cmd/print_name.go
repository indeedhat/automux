package cmd

import (
	"fmt"

	"github.com/indeedhat/automux/internal/config"
)

// PrintNameCommand just prints the session name to std out
func PrintNameCommand(conf *config.Config) error {
	fmt.Println(conf.SessionId)
	return nil
}
