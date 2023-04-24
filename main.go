package main

import (
	"os"

	"github.com/altipla-consulting/errors"
	log "github.com/sirupsen/logrus"

	"github.com/altipla-consulting/ci/internal/commands"
)

func main() {
	if err := commands.CmdRoot.Execute(); err != nil {
		log.Error(err.Error())
		log.Debug(errors.Stack(err))
		os.Exit(1)
	}
}
