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
	var debug, init, detached, printSesionName bool
	var configPath = config.DefaultPath

	flag.BoolVar(&debug, "debug", false, "print tmux commands rather than running them")
	flag.BoolVar(&init, "init", false, "Init the automux config template in the current directory")
	flag.BoolVar(&printSesionName, "print-name", false, "Print the session name if the target directory is a automux directory")
	flag.BoolVar(&detached, "d", false, "Run the automux session detached\nThis will allow you to start an automux session from another session")
	flag.Parse()

	if path := flag.Arg(0); path != "" {
		configPath = path
	}

	if init {
		if err := cmd.InitCmd(); err != nil {
			log.Fatal(err)
		}
		return
	}

	if !config.Exists(configPath) {
		if configPath == config.DefaultPath {
			configPath = config.LegacyPath
			if !config.Exists(configPath) {
				return
			}
		}
		return
	}

	var (
		b bytes.Buffer
		l = log.New(&b, "", 0)
	)

	c, err := config.Load(configPath, l, debug, detached)
	if err != nil {
		log.Fatal("!! invalid automux config !!\n ", err)
	}

	if printSesionName {
		if err := cmd.PrintNameCommand(c); err != nil {
			log.Fatal(err)
		}

		return
	}

	if err := cmd.TriggerCmd(c); err != nil {
		log.Fatal(err)
	}

	fmt.Print(b.String())
}
