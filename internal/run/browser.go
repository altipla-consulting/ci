package run

import (
	"context"
	"os"
	"os/exec"

	"github.com/altipla-consulting/errors"
)

var ErrCannotOpenBrowser = errors.New("run: cannot open browser")

func OpenBrowser(ctx context.Context, url string) error {
	if _, err := exec.LookPath("xdg-open"); err != nil {
		var ee *exec.Error
		if errors.As(err, &ee) && errors.Is(ee.Err, exec.ErrNotFound) {
			return errors.Trace(ErrCannotOpenBrowser)
		}

		return errors.Trace(err)
	}
	cmd := exec.CommandContext(ctx, "xdg-open", url)
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	return errors.Trace(cmd.Start())
}
