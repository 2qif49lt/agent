package middleware

import (
	"net/http"
	"strings"

	"github.com/2qif49lt/agent/api/server/httputils"
	"github.com/2qif49lt/agent/api/types/versions"
	"github.com/2qif49lt/logrus"
	"golang.org/x/net/context"
)

// UserAgentMiddleware is a middleware that
// validates the client user-agent.
type UserAgentMiddleware struct {
	serverVersion string
}

// NewUserAgentMiddleware creates a new UserAgentMiddleware
// with the server version.
func NewUserAgentMiddleware(s string) UserAgentMiddleware {
	return UserAgentMiddleware{
		serverVersion: s,
	}
}

// WrapHandler returns a new handler function wrapping the previous one in the request chain.
func (u UserAgentMiddleware) WrapHandler(handler func(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error) func(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error {
		ctx = context.WithValue(ctx, httputils.UAStringKey, r.Header.Get("User-Agent"))

		if strings.Contains(r.Header.Get("User-Agent"), "Agent-Client/") {
			userAgent := strings.Split(r.Header.Get("User-Agent"), "/")

			// such as Agent-Client/1.7.0 (linux)
			if len(userAgent) == 2 && strings.Contains(userAgent[1], " ") {
				userAgent[1] = strings.Split(userAgent[1], " ")[0]
			}

			if len(userAgent) == 2 && !versions.Equal(u.serverVersion, userAgent[1]) {
				logrus.Debugf("Client and server don't have the same version (client: %s, server: %s)", userAgent[1], u.serverVersion)
			}
		}
		return handler(ctx, w, r, vars)
	}
}
