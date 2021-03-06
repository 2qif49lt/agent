package server

import (
	"github.com/2qif49lt/agent/api/server/httputils"
	"github.com/2qif49lt/agent/api/server/middleware"
	"github.com/2qif49lt/logrus"
)

// handleWithGlobalMiddlwares wraps the handler function for a request with
// the server's global middlewares. The order of the middlewares is backwards,
// meaning that the first in the list will be evaluated last.
func (s *Server) handleWithGlobalMiddlewares(handler httputils.APIFunc) httputils.APIFunc {
	next := handler

	for _, m := range s.middlewares {
		next = m.WrapHandler(next)
	}

	if s.cfg.Logging && logrus.GetLevel() == logrus.DebugLevel {
		next = middleware.DebugRequestMiddleware(next)
	}
	if s.cfg.CertExtenAuth {
		next = middleware.CertExtensionAuthMiddleware(next)
	}

	next = middleware.MissionMiddleware(next)
	return next
}
