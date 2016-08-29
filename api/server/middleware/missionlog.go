package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/2qif49lt/agent/errors"
	//	"github.com/2qif49lt/agent/pkg/eventdb"
	"github.com/2qif49lt/logrus"
	"golang.org/x/net/context"
)

// EventDBMiddleware record  the request mission.
// THIS MIDDLEWARE SHOULD APPEND AT LAST
func EventDBMiddleware(handler func(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error) func(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error {
		logrus.Debugln("EventDBMiddleware enter")
		defer logrus.Debugln("EventDBMiddleware leave")

		paths := strings.Split(r.URL.Path, "/")

		if len(paths) > 0 {
			command := paths[0]
			fmt.Println(command, paths)
		} else {
			return errors.NewErrorWithStatusCode(fmt.Errorf(`url wrong`), http.StatusNotFound)
		}

		return handler(ctx, w, r, vars)
	}
}
