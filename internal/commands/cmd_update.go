package commands

import (
  "github.com/spf13/cobra"
)

var CmdUpdate = &cobra.Command{
  Use:          "update",
  Short:        "Actualiza a la última versión de master borrando todo lo que haya en local",
  RunE: func(cmd *cobra.Command, args []string) error {
    return nil
  },
}
