package commands

import (
  "github.com/spf13/cobra"
)

var CmdPush = &cobra.Command{
  Use:          "push",
  Short:        "Push the commit to production.",
  SilenceUsage: true,
  RunE: func(cmd *cobra.Command, args []string) error {
    return nil
  },
}
