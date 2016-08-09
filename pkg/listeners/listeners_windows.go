package listeners

import (
	"crypto/tls"
	"fmt"
	"net"
	"strings"

	"github.com/2qif49lt/agent/pkg/connections/sockets"
	"github.com/Microsoft/go-winio"
)

// Init creates new listeners for the server.
func Init(proto, addr, socketGroup string, tlsConfig *tls.Config) (net.Listener, error) {

	switch proto {
	case "tcp":
		return sockets.NewTCPSocket(addr, tlsConfig)
	case "npipe":
		// allow Administrators and SYSTEM, plus whatever additional users or groups were specified
		sddl := "D:P(A;;GA;;;BA)(A;;GA;;;SY)"
		if socketGroup != "" {
			for _, g := range strings.Split(socketGroup, ",") {
				sid, err := winio.LookupSidByName(g)
				if err != nil {
					return nil, err
				}
				sddl += fmt.Sprintf("(A;;GRGW;;;%s)", sid)
			}
		}
		c := winio.PipeConfig{
			SecurityDescriptor: sddl,
			MessageMode:        true,  // Use message mode so that CloseWrite() is supported
			InputBufferSize:    65536, // Use 64KB buffers to improve performance
			OutputBufferSize:   65536,
		}
		l, err := winio.ListenPipe(addr, &c)
		if err != nil {
			return nil, err
		}
		return l, err

	default:
		return nil, fmt.Errorf("invalid protocol format: windows only supports tcp and npipe")
	}

	return nil, nil
}
