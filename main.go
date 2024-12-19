package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"

	"github.com/indeedhat/automux/internal/cmd"
	"github.com/indeedhat/automux/internal/config"
)

func main() {
	var configPath = config.DefaultPath

	if path := flag.Arg(0); path != "" {
		configPath = path
	}

	var b bytes.Buffer
	var l = log.New(&b, "", 0)

	root := cmd.Trigger(l, configPath)
	root.AddCommand(cmd.InitC(), cmd.PrintName(l, configPath))

	if err := root.Execute(); err != nil {
		log.Fatal(err)
	}

	fmt.Print(b.String())
}
