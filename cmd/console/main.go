package main

import (
	"os"

	"github.com/dannyexmakina/go-stencil-laravalian/app/command"
)

func main() {
	if err := command.RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
