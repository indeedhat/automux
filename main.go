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
	var debug, init bool

	flag.BoolVar(&debug, "debug", false, "print tmux commands rather than running them")
	flag.BoolVar(&init, "init", false, "Init the automux config template in the current directory")
	flag.Parse()

	if init {
		if err := cmd.InitCmd(); err != nil {
			log.Fatal(err)
		}
		return
	}

	if !config.Exists() {
		return
	}

	var (
		b bytes.Buffer
		l = log.New(&b, "", 0)
	)

	c, err := config.Load(config.DefaultPath, l, debug)
	if err != nil {
		log.Fatal("!! invalid automux config !!\n ", err)
	}

	if err := cmd.TriggerCmd(c); err != nil {
		log.Fatal(err)
	}

	fmt.Print(b.String())
}
