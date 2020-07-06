package query

import (
	"strings"

	"libs.altipla.consulting/errors"

	"github.com/altipla-consulting/ci/internal/run"
)

var org, repo string

func IsGerrit() (bool, error) {
	remote, err := run.GitCaptureOutput("remote", "get-url", "origin")
	if err != nil {
		return false, errors.Trace(err)
	}
	return strings.Contains(remote, "gerrit.altiplaconsulting.net"), nil
}

func extractGitHub() error {
	if org != "" {
		return nil
	}

	remote, err := run.GitCaptureOutput("remote", "get-url", "origin")
	if err != nil {
		return errors.Trace(err)
	}

	parts := strings.Split(remote, "/")
	org = parts[0][len("git@github.com:"):]
	repo = parts[1][:len(parts[1])-len(".git\n")]
	return nil
}

func CurrentOrg() (string, error) {
	if err := extractGitHub(); err != nil {
		return "", errors.Trace(err)
	}
	return org, nil
}

func CurrentRepo() (string, error) {
	if err := extractGitHub(); err != nil {
		return "", errors.Trace(err)
	}
	return repo, nil
}

func CurrentBranch() (string, error) {
	branch, err := run.GitCaptureOutput("rev-parse", "--abbrev-ref", "HEAD")
	if err != nil {
		return "", errors.Trace(err)
	}
	return branch, nil
}
