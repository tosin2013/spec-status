package main

import (
	"os"

	"github.com/tosin2013/spec-status/internal/cli"
)

func main() {
	os.Exit(cli.Run(os.Args[1:], "", os.Stdout, os.Stderr))
}
