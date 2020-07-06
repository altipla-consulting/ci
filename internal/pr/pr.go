package pr

import (
	"context"

	"github.com/google/go-github/v32/github"
	"golang.org/x/oauth2"
	"libs.altipla.consulting/errors"

	"github.com/altipla-consulting/ci/internal/login"
	"github.com/altipla-consulting/ci/internal/query"
)

var client *github.Client

func initClient(ctx context.Context) error {
	if client != nil {
		return nil
	}

	auth, err := login.ReadAuthConfig()
	if err != nil {
		return errors.Trace(err)
	}
	if auth == nil {
		return errors.Errorf("Inicia sesión con `ci login` antes de interactuar con GitHub")
	}
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: auth.AccessToken})
	client = github.NewClient(oauth2.NewClient(ctx, ts))

	return nil
}

func ListBranches(ctx context.Context) ([]string, error) {
	if err := initClient(ctx); err != nil {
		return nil, errors.Trace(err)
	}

	org, err := query.CurrentOrg()
	if err != nil {
		return nil, errors.Trace(err)
	}
	repo, err := query.CurrentRepo()
	if err != nil {
		return nil, errors.Trace(err)
	}
	prs, _, err := client.PullRequests.List(ctx, org, repo, nil)
	if err != nil {
		return nil, errors.Trace(err)
	}

	var names []string
	for _, pr := range prs {
		names = append(names, pr.GetHead().GetRef())
	}
	return names, nil
}
