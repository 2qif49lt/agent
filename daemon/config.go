package daemon

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"strings"
	"sync"

	"github.com/2qif49lt/agent/api/types"
	"github.com/2qif49lt/logrus"
	flag "github.com/2qif49lt/pflag"
)

// Config defines the configuration of a agent daemon.
type Config struct {
	AgentID           string
	CorsHeaders       string
	NoTLSClientVerify bool // 是否要求客户端验证
	SocketGroup       string
	reloadLock        sync.Mutex
}

// ReloadConfiguration reads the configuration in the host and reloads the daemon and server.
func ReloadConfiguration(configFile string, flags *flag.FlagSet, reload func(*Config)) error {
	logrus.Infof("Got signal to reload configuration, reloading from: %s", configFile)

	return nil
}

// InstallFlags adds command-line options to the top-level flag parser for
// the current process.
// Subsequent calls to `flag.Parse` will populate config with values parsed
// from the command-line.
func (config *Config) InstallFlags(cmd *flag.FlagSet) {
	cmd.StringVar(&config.AgentID, "agent-id", "", "Set CORS headers in the remote API")
	cmd.StringVar(&config.CorsHeaders, "api-cors-header", "", "Set CORS headers in the remote API")
	cmd.BoolVar(&config.NoTLSClientVerify, "noverify", false, "DO NOT verify client certificate")
	cmd.StringVar(&config.SocketGroup, "group", "agentd", "Group name for the unix socket")
}
