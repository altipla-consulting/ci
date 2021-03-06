package commands

import (
	"context"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"libs.altipla.consulting/collections"
	"libs.altipla.consulting/datetime"
	"libs.altipla.consulting/errors"

	"github.com/altipla-consulting/ci/internal/login"
	"github.com/altipla-consulting/ci/internal/pr"
	"github.com/altipla-consulting/ci/internal/prompt"
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

		mainBranch, err := query.MainBranch()
		if err != nil {
			return errors.Trace(err)
		}
		branch, err := query.CurrentBranch()
		if err != nil {
			return errors.Trace(err)
		}
		if branch != mainBranch {
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
		} else {
			auth, err := login.ReadAuthConfig()
			if err != nil {
				return errors.Trace(err)
			}
			if auth == nil {
				return errors.Errorf("Inicia sesión con `ci login` antes de interactuar con GitHub")
			}
			t := time.Now().In(datetime.EuropeMadrid()).Format("0405")
			branch = fmt.Sprintf("f/%s-%s", auth.Username, t)
			if err := run.Git("checkout", "-b", branch); err != nil {
				return errors.Trace(err)
			}
		}

		if err := run.Git("push", "-u", "origin", branch); err != nil {
			return errors.Trace(err)
		}

		last, err := query.LastCommitMessage()
		if err != nil {
			return errors.Trace(err)
		}

		fmt.Println()
		title, err := prompt.TextDefault("Título del PR", last)
		if err != nil {
			return errors.Trace(err)
		}
		link, err := pr.Create(ctx, title)
		if err != nil {
			return errors.Trace(err)
		}

		log.Info()
		log.Info("Se ha creado un nuevo PR en el repo de GitHub.")
		log.Info("\t" + link)
		log.Info()

		if err := run.OpenBrowser(link); err != nil {
			return errors.Trace(err)
		}

		return nil
	},
}
