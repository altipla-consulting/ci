package commands

import (
	"context"

	"github.com/spf13/cobra"
	"libs.altipla.consulting/collections"
	"libs.altipla.consulting/errors"

	"github.com/altipla-consulting/ci/internal/pr"
	"github.com/altipla-consulting/ci/internal/query"
	"github.com/altipla-consulting/ci/internal/run"
)

var CmdPush = &cobra.Command{
	Use:   "push",
	Short: "Envía el commit a Gerrit/GitHub o crea un nuevo PR si no existe",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()

		gerrit, err := query.IsGerrit()
		if err != nil {
			return errors.Trace(err)
		}

		if gerrit {
			if err := run.Git("push", "origin", "HEAD:refs/for/master"); err != nil {
				return errors.Trace(err)
			}
			return nil
		}

		branch, err := query.CurrentBranch()
		if err != nil {
			return errors.Trace(err)
		}
		if branch != "master" {
			branches, err := pr.ListBranches(ctx)
			if err != nil {
				return errors.Trace(err)
			}
			if collections.HasString(branches, branch) {
				// La rama tiene un PR abierto, enviamos el nuevo commit que automáticamente
				// sale en la interfaz de PRs.
				if err := run.Git("push"); err != nil {
					return errors.Trace(err)
				}
				return nil
			}
		}

		return nil
	},
}
