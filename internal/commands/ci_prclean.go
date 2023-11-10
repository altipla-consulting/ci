package commands

import (
	"strings"

	"github.com/altipla-consulting/errors"
	"github.com/spf13/cobra"
	"golang.org/x/exp/slices"

	"github.com/altipla-consulting/ci/internal/query"
	"github.com/altipla-consulting/ci/internal/run"
	log "github.com/sirupsen/logrus"
)

var cmdPRClean = &cobra.Command{
	Use:     "prclean",
	Aliases: []string{"prc"},
	Short:   "Clean all local branches that have no equivalent in the remote.",
	Example: "ci prclean",
	RunE: func(cmd *cobra.Command, args []string) error {
		gerrit, err := query.IsGerrit()
		if err != nil {
			return errors.Trace(err)
		}
		if gerrit {
			return errors.Errorf("Gerrit does not use PRs")
		}

		if err := run.GitContext(cmd.Context(), "fetch"); err != nil {
			return errors.Trace(err)
		}
		if err := run.GitContext(cmd.Context(), "remote", "prune", "origin"); err != nil {
			return errors.Trace(err)
		}

		current, err := query.CurrentBranch()
		if err != nil {
			return errors.Trace(err)
		}
		local, err := query.LocalBranches(cmd.Context())
		if err != nil {
			return errors.Trace(err)
		}
		remote, err := query.RemoteBranches(cmd.Context())
		if err != nil {
			return errors.Trace(err)
		}
		var clean int
		var keep []string
		for _, branch := range local {
			if slices.Contains(remote, branch) {
				keep = append(keep, branch)
				continue
			}

			clean++
			if branch == current {
				log.Info("Changing branch to the main one to clean the old branch")
				main, err := query.MainBranch()
				if err != nil {
					return errors.Trace(err)
				}
				if err := run.GitContext(cmd.Context(), "checkout", main); err != nil {
					return errors.Trace(err)
				}
			}

			if err := run.GitContext(cmd.Context(), "branch", "-D", branch); err != nil {
				return errors.Trace(err)
			}
		}

		log.Infof("Cleaned %d branches, %d branches kept (%s)", clean, len(keep), strings.Join(keep, ", "))

		return nil
	},
}
