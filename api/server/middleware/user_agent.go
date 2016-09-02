package middleware

import (
	"context"
	"net/http"

	"github.com/2qif49lt/agent/api/server/httputils"
	"github.com/2qif49lt/logrus"
)

// UserAgentMiddleware is a middleware that
// validates the client user-agent.
type UserAgentMiddleware struct {
	// count
	c int
}

// NewUserAgentMiddleware creates a new UserAgentMiddleware
// with the server version.
func NewUserAgentMiddleware() *UserAgentMiddleware {
	return &UserAgentMiddleware{}
}

// WrapHandler returns a new handler function wrapping the previous one in the request chain.
func (u *UserAgentMiddleware) WrapHandler(handler func(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error) func(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error {
		logrus.Debugln("UserAgentMiddleware enter")
		defer logrus.Debugln("UserAgentMiddleware leave")

		ctx = context.WithValue(ctx, httputils.UAStringKey, r.Header.Get("User-Agent"))

		return handler(ctx, w, r, vars)
	}
}
