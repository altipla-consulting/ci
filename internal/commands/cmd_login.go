package commands

import (
	"github.com/altipla-consulting/errors"
	"github.com/spf13/cobra"

	"github.com/altipla-consulting/ci/internal/login"
)

var cmdLogin = &cobra.Command{
	Use:     "login",
	Short:   "Login to GitHub for all operations related to that type of repos.",
	Example: "ci login",
	RunE: func(cmd *cobra.Command, args []string) error {
		return errors.Trace(login.Start(cmd.Context()))
	},
}
