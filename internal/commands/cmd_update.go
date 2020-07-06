package commands

import (
	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
	"libs.altipla.consulting/errors"

	"github.com/altipla-consulting/ci/internal/run"
)

var CmdUpdate = &cobra.Command{
	Use:   "update",
	Short: "Actualiza a la última versión de master borrando todo lo que haya en local",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := run.Git("fetch", "origin"); err != nil {
			return errors.Trace(err)
		}

		status, err := run.GitCaptureOutput("status", "-s")
		if err != nil {
			return errors.Trace(err)
		}
		if len(status) > 0 {
			var keep bool
			prompt := &survey.Confirm{
				Message: "El proyecto tiene cambios. ¿Estás seguro de que deseas borrar todo y pasar a master?",
			}
			if err := survey.AskOne(prompt, &keep); err != nil {
				return errors.Trace(err)
			}
			if !keep {
				return nil
			}
		}

		if err := run.Git("checkout", "master"); err != nil {
			return errors.Trace(err)
		}
		if err := run.Git("reset", "--hard", "origin/master"); err != nil {
			return errors.Trace(err)
		}
		if err := run.Git("stash", "-q", "--include-untracked"); err != nil {
			return errors.Trace(err)
		}
		if err := run.Git("stash", "drop", "-q"); err != nil {
			return errors.Trace(err)
		}

		return nil
	},
}
