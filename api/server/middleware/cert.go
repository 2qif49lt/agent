package middleware

import (
	"context"
	"encoding/asn1"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/2qif49lt/agent/api/server/httputils"
	"github.com/2qif49lt/agent/errors"
	"github.com/2qif49lt/logrus"
)

const (
	defaultExtenAuth = "ping info"
	defaultExtenOID  = "1.2.3.4"
)

// CertExtensionAuthMiddleware check client's certificate's custom field 1.2.3.4 for authenticating.
// 可能考虑只对插件进行类似的授权鉴定
func CertExtensionAuthMiddleware(handler func(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error) func(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error {
		logrus.Debugln("CertExtensionAuthMiddleware enter")
		defer logrus.Debugln("CertExtensionAuthMiddleware leave")

		if r.TLS != nil {
			if cert := r.TLS.PeerCertificates[0]; cert != nil {
				bauth := false
				auth := ""

				for _, exten := range cert.Extensions {
					if exten.Id.String() == defaultExtenOID {
						_, err := asn1.Unmarshal(exten.Value, &auth)
						if err != nil {
							return errors.NewErrorWithStatusCode(err, http.StatusInternalServerError)
						}
						break
					}
				}
				logrus.Debugln("CertExtensionAuthMiddleware auth:", auth)

				if auth == "" {
					auth = defaultExtenAuth
				}

				command := httputils.CommandFromRequest(r)

				if len(command) > 0 {
					authexps := strings.Split(auth, " ")

					for _, authexp := range authexps {
						authexp = fmt.Sprintf(`^%s$`, authexp)
						if match, err := regexp.MatchString(authexp, command); match == true && err == nil {
							bauth = true
							break
						}
					}

				} else {
					bauth = true
				}

				if bauth == false {
					return errors.NewErrorWithStatusCode(fmt.Errorf(`certificate 's extension authenticate fail, auth:%s`, auth),
						http.StatusUnauthorized)
				}

			}

		}
		return handler(ctx, w, r, vars)
	}
}
