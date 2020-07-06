package commands

import (
  "github.com/spf13/cobra"
)

var CmdPush = &cobra.Command{
  Use:          "push",
  Short:        "Envía el commit a Gerrit/GitHub o crea un nuevo PR si no existe",
  RunE: func(cmd *cobra.Command, args []string) error {
    return nil
  },
}
