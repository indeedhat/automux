package main

import (
	"bytes"
	"fmt"
	"log"

	"github.com/indeedhat/automux/internal/cmd"
)

func main() {
	var b bytes.Buffer
	var l = log.New(&b, "", 0)

	root := cmd.Trigger(l)
	root.AddCommand(cmd.Init(), cmd.PrintName(l))

	if err := root.Execute(); err != nil {
		log.Fatal(err)
	}

	fmt.Print(b.String())
}
