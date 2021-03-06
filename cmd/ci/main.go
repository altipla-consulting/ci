package main

import (
	"os"

	log "github.com/sirupsen/logrus"
	"libs.altipla.consulting/errors"

	"github.com/altipla-consulting/ci/internal/commands"
)

func main() {
	if err := commands.CmdRoot.Execute(); err != nil {
		log.Error(errors.Stack(err))
		os.Exit(1)
	}
}
