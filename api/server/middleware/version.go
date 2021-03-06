package middleware

import (
	"context"
	"fmt"
	"net/http"
	"runtime"

	"github.com/2qif49lt/agent/api"
	"github.com/2qif49lt/agent/api/server/httputils"
	"github.com/2qif49lt/agent/api/types/versions"
	"github.com/2qif49lt/logrus"
)

type badRequestError struct {
	error
}

func (badRequestError) HTTPErrorStatusCode() int {
	return http.StatusBadRequest
}

// VersionMiddleware is a middleware that
// validates the client and server versions.
// defaultVersion: defautl version if client dont take，默认对方提供服务的版本
// minVersion: minimum api version server support
// serverVersion: the server api's version,serverVersion可能大于defaultVersion,如测试新功能。
type VersionMiddleware struct {
	serverVersion  string
	defaultVersion string
	minVersion     string
}

// NewVersionMiddleware creates a new VersionMiddleware
// with the default versions.
func NewVersionMiddleware(s, d, m string) VersionMiddleware {
	return VersionMiddleware{
		serverVersion:  s,
		defaultVersion: d,
		minVersion:     m,
	}
}

// WrapHandler returns a new handler function wrapping the previous one in the request chain.
func (v VersionMiddleware) WrapHandler(handler func(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error) func(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error {
		logrus.Debugln("VersionMiddleware enter")
		defer logrus.Debugln("VersionMiddleware leave")

		apiVersion := r.Header.Get(api.HEADER_VERISON)
		if apiVersion == "" {
			logrus.WithField("remote", r.RemoteAddr).Warnln("api version is empty")

			apiVersion = v.defaultVersion
		}

		header := fmt.Sprintf("Agentd/%s (%s)", v.defaultVersion, runtime.GOOS)
		w.Header().Set("Server", header)

		if versions.GreaterThan(apiVersion, v.serverVersion) {
			return badRequestError{fmt.Errorf("client is newer than server (client API version: %s, server API version: %s)", apiVersion, v.serverVersion)}
		}
		if versions.LessThan(apiVersion, v.minVersion) {
			return badRequestError{fmt.Errorf("client version %s is too old. Minimum supported API version is %s, please upgrade your client to a newer version", apiVersion, v.minVersion)}
		}

		httputils.Put(vars, httputils.CLI_API_VERSION, apiVersion)

		return handler(ctx, w, r, vars)
	}

}
