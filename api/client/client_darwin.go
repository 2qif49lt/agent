package client

import (
	"fmt"
	"github.com/2qif49lt/agent/cfg"
)

// DefaultAgentdHost
const DefaultAgentdHost = fmt.Sprintf(`tcp://127.0.0.1:%d`, cfg.DefaultAgentdListenPort)
