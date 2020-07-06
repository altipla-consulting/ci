package commands

import (
  "github.com/spf13/cobra"
)

var CmdUpdate = &cobra.Command{
  Use:          "update",
  Short:        "Update to the latest master version erasing everything",
  RunE: func(cmd *cobra.Command, args []string) error {
    return nil
  },
}
