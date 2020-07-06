package commands

import (
  "github.com/spf13/cobra"
)

var CmdDraft = &cobra.Command{
  Use:          "draft",
  Short:        "Crea un nuevo commit WIP (Work-In-Progress) en caso de Gerrit o un PR draft en caso de GitHub",
  RunE: func(cmd *cobra.Command, args []string) error {
    return nil
  },
}
