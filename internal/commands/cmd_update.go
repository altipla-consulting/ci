package commands

import (
	"fmt"

	"github.com/altipla-consulting/errors"
	"github.com/spf13/cobra"

	"github.com/altipla-consulting/ci/internal/prompt"
	"github.com/altipla-consulting/ci/internal/query"
	"github.com/altipla-consulting/ci/internal/run"
)

var flagForce bool

func init() {
	cmdUpdate.Flags().BoolVarP(&flagForce, "force", "f", false, "Fuerza la actualización aunque haya cambios pendientes. WARNING: Es una operación destructiva.")
}

var cmdUpdate = &cobra.Command{
	Use:   "update",
	Short: "Actualiza a la última versión de master borrando todo lo que haya en local",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := run.GitContext(cmd.Context(), "fetch", "origin"); err != nil {
			return errors.Trace(err)
		}

		mainBranch, err := query.MainBranch()
		if err != nil {
			return errors.Trace(err)
		}

		status, err := run.GitCaptureOutputContext(cmd.Context(), "status", "-s")
		if err != nil {
			return errors.Trace(err)
		}
		if len(status) > 0 && !flagForce {
			keep, err := prompt.Confirm(fmt.Sprintf("El proyecto tiene cambios. ¿Estás seguro de que deseas borrar todo y pasar a %s?", mainBranch))
			if err != nil {
				return errors.Trace(err)
			}
			if !keep {
				return nil
			}
		}

		if err := run.GitContext(cmd.Context(), "checkout", "--", "."); err != nil {
			return errors.Trace(err)
		}
		if err := run.GitContext(cmd.Context(), "checkout", mainBranch); err != nil {
			return errors.Trace(err)
		}
		if err := run.GitContext(cmd.Context(), "reset", "--hard", "origin/"+mainBranch); err != nil {
			return errors.Trace(err)
		}

		// Remove any untracked file remaining in the working directory
		status, err = run.GitCaptureOutputContext(cmd.Context(), "status", "-s")
		if err != nil {
			return errors.Trace(err)
		}
		if len(status) > 0 {
			if err := run.GitContext(cmd.Context(), "stash", "-q", "--include-untracked"); err != nil {
				return errors.Trace(err)
			}
			if err := run.GitContext(cmd.Context(), "stash", "drop", "-q"); err != nil {
				return errors.Trace(err)
			}
		}

		return nil
	},
}
