package login

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/altipla-consulting/errors"
	"github.com/google/go-github/v43/github"
	"github.com/segmentio/ksuid"
	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	oauth2github "golang.org/x/oauth2/github"

	"github.com/altipla-consulting/ci/internal/prompt"
	"github.com/altipla-consulting/ci/internal/run"
)

func Start(ctx context.Context) error {
	auth, err := ReadAuthConfig()
	if err != nil {
		return errors.Trace(err)
	}

	if auth != nil {
		keep, err := prompt.Confirm(fmt.Sprintf("Ya tienes una sesión configurada con %s. ¿Estás seguro de que deseas iniciar sesión de nuevo?", auth.Username))
		if err != nil {
			return errors.Trace(err)
		}
		if !keep {
			return nil
		}
	}

	conf := &oauth2.Config{
		ClientID:     "d7d38c193ed7aeebf5c2",
		ClientSecret: "b6e6b7939c67e4fdf9803563e1194b29e47b3c62",
		Scopes:       []string{"repo"},
		Endpoint:     oauth2github.Endpoint,
	}

	var server *http.Server
	handler := func(w http.ResponseWriter, r *http.Request) error {
		fmt.Fprintln(w, "Cierra esta pestaña. La autenticación ya ha terminado.")
		go func() {
			time.Sleep(2 * time.Second)
			_ = server.Close()
		}()

		token, err := conf.Exchange(ctx, r.FormValue("code"))
		if err != nil {
			return errors.Trace(err)
		}
		ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token.AccessToken})
		client := github.NewClient(oauth2.NewClient(ctx, ts))

		user, _, err := client.Users.Get(ctx, "")
		if err != nil {
			return errors.Trace(err)
		}

		auth := AuthConfig{
			AccessToken: token.AccessToken,
			Username:    user.GetLogin(),
		}
		if err := auth.Save(); err != nil {
			return errors.Trace(err)
		}

		log.Info()
		log.Info("Autorización recibida. Sesión iniciada correctamente.")
		log.Info()

		return nil
	}

	errch := make(chan error, 1)
	server = &http.Server{
		Addr: ":45678",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			errch <- handler(w, r)
		}),
	}

	state := ksuid.New().String()
	url := conf.AuthCodeURL(state, oauth2.AccessTypeOffline)
	log.Info("Abriendo navegador para iniciar sesión...")
	log.Info(url)
	if err := run.OpenBrowser(ctx, url); err != nil {
		return errors.Trace(err)
	}

	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return errors.Trace(err)
	}
	return errors.Trace(<-errch)
}
