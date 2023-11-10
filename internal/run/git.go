package run

import (
	"context"
	"os"
	"os/exec"
	"strings"

	"github.com/altipla-consulting/errors"
	log "github.com/sirupsen/logrus"
)

func Git(ctx context.Context, args ...string) error {
	log.Debugf("RUN: git %s", strings.Join(args, " "))

	cmd := exec.CommandContext(ctx, "git", args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return errors.Trace(cmd.Run())
}

func GitCaptureOutput(ctx context.Context, args ...string) (string, error) {
	log.Debugf("RUN: git %s", strings.Join(args, " "))

	cmd := exec.CommandContext(ctx, "git", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", errors.Trace(err)
	}

	return strings.TrimSpace(string(output)), nil
}
