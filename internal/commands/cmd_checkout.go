package commands

import (
  "github.com/spf13/cobra"
)

var CmdCheckout = &cobra.Command{
  Use:          "checkout",
  Short:        "Checkout the code to a remote GitHub pull request",
  RunE: func(cmd *cobra.Command, args []string) error {
    return nil
  },
}
