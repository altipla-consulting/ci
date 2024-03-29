package query

import (
	"context"
	"net/url"
	"os/exec"
	"path"
	"regexp"
	"strings"

	"github.com/altipla-consulting/errors"

	"github.com/altipla-consulting/ci/internal/run"
)

var org, repo string

func IsGerrit(ctx context.Context) (bool, error) {
	remote, err := run.GitCaptureOutput(ctx, "remote", "get-url", "origin")
	if err != nil {
		return false, errors.Trace(err)
	}
	return strings.Contains(remote, "gerrit.altipla.consulting"), nil
}

func MainBranch(ctx context.Context) (string, error) {
	branch, err := run.GitCaptureOutput(ctx, "branch", "-a")
	if err != nil {
		return "", errors.Trace(err)
	}
	mainBranch := "master"
	if strings.Contains(branch, "remotes/origin/main") {
		mainBranch = "main"
	}
	return mainBranch, nil
}

var scpSyntaxRe = regexp.MustCompile(`^([a-zA-Z0-9_]+)@([a-zA-Z0-9._-]+):(.*)$`)

func extractGitHub(ctx context.Context) error {
	if org != "" {
		return nil
	}

	remote, err := run.GitCaptureOutput(ctx, "remote", "get-url", "origin")
	if err != nil {
		return errors.Trace(err)
	}
	var remoteURL *url.URL
	if m := scpSyntaxRe.FindStringSubmatch(remote); m != nil {
		// Match SCP-like syntax and convert it to a URL.
		// Eg, "git@github.com:user/repo" becomes
		// "ssh://git@github.com/user/repo".
		remoteURL = &url.URL{
			Scheme: "ssh",
			User:   url.User(m[1]),
			Host:   m[2],
			Path:   m[3],
		}
	} else {
		remoteURL, err = url.Parse(remote)
		if err != nil {
			return errors.Trace(err)
		}
	}

	parts := strings.SplitN(strings.TrimSuffix(remoteURL.Path, ".git"), "/", 2)
	org = parts[0]
	repo = parts[1]
	return nil
}

func CurrentOrg(ctx context.Context) (string, error) {
	if err := extractGitHub(ctx); err != nil {
		return "", errors.Trace(err)
	}
	return org, nil
}

func CurrentRepo(ctx context.Context) (string, error) {
	if err := extractGitHub(ctx); err != nil {
		return "", errors.Trace(err)
	}
	return repo, nil
}

func CurrentBranch(ctx context.Context) (string, error) {
	branch, err := run.GitCaptureOutput(ctx, "rev-parse", "--abbrev-ref", "HEAD")
	if err != nil {
		return "", errors.Trace(err)
	}
	return branch, nil
}

func LastCommitMessage(ctx context.Context) (string, error) {
	msg, err := run.GitCaptureOutput(ctx, "log", "-1", "--pretty=%B")
	if err != nil {
		return "", errors.Trace(err)
	}
	return msg, nil
}

func BranchExists(ctx context.Context, name string) (bool, error) {
	if err := run.Git(ctx, "show-ref", "--verify", "--quiet", "refs/heads/"+name); err != nil {
		var exit *exec.ExitError
		if errors.As(err, &exit) {
			if exit.ProcessState.ExitCode() == 1 {
				return false, nil
			}
		}
		return false, errors.Trace(err)
	}
	return true, nil
}

func LocalBranches(ctx context.Context) ([]string, error) {
	branches, err := run.GitCaptureOutput(ctx, "branch", "--format=%(refname:short)")
	if err != nil {
		return nil, errors.Trace(err)
	}
	return strings.Split(strings.TrimSpace(branches), "\n"), nil
}

func RemoteBranches(ctx context.Context) ([]string, error) {
	branches, err := run.GitCaptureOutput(ctx, "branch", "-r", "--format=%(refname:short)")
	if err != nil {
		return nil, errors.Trace(err)
	}
	names := strings.Split(strings.TrimSpace(branches), "\n")
	for i, name := range names {
		names[i] = path.Base(name)
	}
	return names, nil
}
