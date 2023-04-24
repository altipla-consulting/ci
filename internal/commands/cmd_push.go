package commands

import (
	"github.com/altipla-consulting/errors"
	"github.com/spf13/cobra"

	"github.com/altipla-consulting/ci/internal/query"
	"github.com/altipla-consulting/ci/internal/run"
)

var cmdPush = &cobra.Command{
	Use:     "push",
	Short:   "Envía el commit a Gerrit/GitHub.",
	Example: "ci push",
	RunE: func(cmd *cobra.Command, args []string) error {
		gerrit, err := query.IsGerrit()
		if err != nil {
			return errors.Trace(err)
		}
		mainBranch, err := query.MainBranch()
		if err != nil {
			return errors.Trace(err)
		}

		if gerrit {
			if err := run.Git("push", "origin", "HEAD:refs/for/"+mainBranch); err != nil {
				return errors.Trace(err)
			}
			return nil
		}

		return errors.Trace(run.Git("push"))
	},
}
