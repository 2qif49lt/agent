// +build !windows

package listeners

import (
	"crypto/tls"
	"fmt"
	"net"

	"github.com/2qif49lt/agent/pkg/connections/sockets"
)

// Init creates new listeners for the server.
// TODO: Clean up the fact that socketGroup and tlsConfig aren't always used.
func Init(proto, addr, socketGroup string, tlsConfig *tls.Config) (net.Listener, error) {
	switch proto {
	case "tcp":
		return sockets.NewTCPSocket(addr, tlsConfig)
	case "unix":
		l, err := sockets.NewUnixSocket(addr, socketGroup)
		if err != nil {
			return nil, fmt.Errorf("can't create unix socket %s: %v", addr, err)
		}
		return l, err
	default:

		return nil, fmt.Errorf("invalid protocol format: %q", proto)
	}

	return nil, nil
}
