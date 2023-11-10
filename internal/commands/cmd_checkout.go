package commands

import (
	"fmt"
	"strconv"

	"github.com/altipla-consulting/errors"
	"github.com/spf13/cobra"

	"github.com/altipla-consulting/ci/internal/pr"
	"github.com/altipla-consulting/ci/internal/query"
	"github.com/altipla-consulting/ci/internal/run"
)

var cmdCheckout = &cobra.Command{
	Use:     "checkout",
	Short:   "Establece el código a la versión exacta de un Pull Request en GitHub",
	Example: "ci checkout 123",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return errors.Errorf("Especifica como primer argumento el ID del PR que quieres descargar: %v: %v", args[0], err)
		}

		branch, err := pr.Branch(cmd.Context(), id)
		if err != nil {
			return errors.Trace(err)
		}

		exists, err := query.BranchExists(branch)
		if err != nil {
			return errors.Trace(err)
		}
		if exists {
			if err := run.GitContext(cmd.Context(), "branch", "-D", branch); err != nil {
				return errors.Trace(err)
			}
		}
		if err := run.GitContext(cmd.Context(), "fetch", "origin", fmt.Sprintf("pull/%d/head:%s", id, branch)); err != nil {
			return errors.Trace(err)
		}

		if err := run.GitContext(cmd.Context(), "checkout", branch); err != nil {
			return errors.Trace(err)
		}

		return nil
	},
}

var cmdCheckoutShort = &cobra.Command{
	Use:     "co",
	Short:   cmdCheckout.Short,
	Example: "ci co 123",
	Args:    cmdCheckout.Args,
	RunE:    cmdCheckout.RunE,
}
