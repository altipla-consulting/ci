package commands

import (
	"github.com/spf13/cobra"
)

var CmdCheckout = &cobra.Command{
	Use:   "checkout",
	Short: "Establece el código a la versión exacta de un Pull Request en GitHub",
	RunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}
