package main

import (
	"os"
	
	"github.com/altipla-consulting/ci/internal/commands"
)

func main() {
	if err := commands.CmdRoot.Execute(); err != nil {
		os.Exit(1)
	}
}
