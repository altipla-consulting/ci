package commands

import (
	"fmt"
	"strings"
	"time"

	"github.com/altipla-consulting/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"golang.org/x/exp/slices"

	"github.com/altipla-consulting/ci/internal/login"
	"github.com/altipla-consulting/ci/internal/pr"
	"github.com/altipla-consulting/ci/internal/query"
	"github.com/altipla-consulting/ci/internal/run"
)

var europeMadrid *time.Location

func init() {
	var err error
	europeMadrid, err = time.LoadLocation("Europe/Madrid")
	if err != nil {
		panic(fmt.Sprintf("cannot load location Europe/Madrid: %v", err))
	}
}

var cmdPR = &cobra.Command{
	Use:     "pr",
	Short:   "Create a new branch and send a PR to the main branch.",
	Example: "ci pr",
	RunE: func(cmd *cobra.Command, args []string) error {
		gerrit, err := query.IsGerrit(cmd.Context())
		if err != nil {
			return errors.Trace(err)
		}
		if gerrit {
			return errors.Errorf("Gerrit does not use PRs")
		}

		mainBranch, err := query.MainBranch(cmd.Context())
		if err != nil {
			return errors.Trace(err)
		}
		branch, err := query.CurrentBranch(cmd.Context())
		if err != nil {
			return errors.Trace(err)
		}
		if branch != mainBranch {
			branches, err := pr.ListBranches(cmd.Context())
			if err != nil {
				return errors.Trace(err)
			}
			if slices.Contains(branches, branch) {
				// La rama tiene un PR abierto, enviamos el nuevo commit que automáticamente
				// sale en la interfaz de PRs.
				if err := run.Git(cmd.Context(), "push"); err != nil {
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
				return errors.Errorf("Inicia sesión con `ci login` antes de interactuar con GitContextHcmd.Context(), ub")
			}
			t := time.Now().In(europeMadrid).Format("0405")
			branch = fmt.Sprintf("f/%s-%s", auth.Username, t)
			if err := run.Git(cmd.Context(), "checkout", "-b", branch); err != nil {
				return errors.Trace(err)
			}
		}

		if err := run.Git(cmd.Context(), "push", "-u", "origin", branch); err != nil {
			return errors.Trace(err)
		}

		last, err := query.LastCommitMessage(cmd.Context())
		if err != nil {
			return errors.Trace(err)
		}
		parts := strings.SplitN(last, "\n\n", 2)

		var body string
		if len(parts) > 1 {
			body = parts[1]
		}
		link, err := pr.Create(cmd.Context(), parts[0], body)
		if err != nil {
			return errors.Trace(err)
		}

		log.Info()
		log.Info("Se ha creado un nuevo PR en el repo de GitHub.")
		log.Info("\t" + link)
		log.Info()

		if err := run.OpenBrowser(cmd.Context(), link); err != nil && !errors.Is(err, run.ErrCannotOpenBrowser) {
			return errors.Trace(err)
		}

		return nil
	},
}
