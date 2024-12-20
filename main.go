package main

import (
	"bytes"
	"context"
	"fmt"
	"log"

	"github.com/indeedhat/automux/internal/cmd"
)

func main() {
	var b bytes.Buffer
	var l = log.New(&b, "", 0)

	ctx := context.WithValue(context.Background(), "logger", l)

	root := cmd.Trigger()
	root.AddCommand(cmd.Init(), cmd.PrintName())

	if err := root.ExecuteContext(ctx); err != nil {
		log.Fatal(err)
	}

	fmt.Print(b.String())
}
