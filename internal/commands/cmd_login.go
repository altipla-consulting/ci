package commands

import (
	"context"

	"github.com/spf13/cobra"
	"libs.altipla.consulting/errors"

	"github.com/altipla-consulting/ci/internal/login"
)

var CmdLogin = &cobra.Command{
	Use:   "login",
	Short: "Inicia sesión global en GitHub para todas las operaciones relacionadas con ese tipo de repos",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()

		if err := login.Start(ctx); err != nil {
			return errors.Trace(err)
		}

		return nil
	},
}
