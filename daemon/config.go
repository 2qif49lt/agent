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
	AgentID string

	Pidfile string

	EnableCors  bool
	CorsHeaders string

	reloadLock sync.Mutex
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
	cmd.StringVar(&config.AgentID, "agentid", "", "Set CORS headers in the remote API")
	cmd.StringVar(&config.Pidfile, "pidfile", defaultPidFile, "Path to use for daemon PID file")
	cmd.StringVar(&config.CorsHeaders, "api-cors-header", "", "Set CORS headers in the remote API")
}
