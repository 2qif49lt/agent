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

const (
	defaultNetworkMtu = 1500
)

// CommonTLSOptions defines TLS configuration for the daemon server.
// It includes json tags to deserialize configuration from a file
// using the same names that the flags in the command line use.
type CommonTLSOptions struct {
	CAFile   string `json:"tlscacert,omitempty"`
	CertFile string `json:"tlscert,omitempty"`
	KeyFile  string `json:"tlskey,omitempty"`
}

// CommonConfig defines the configuration of a agent daemon which is
// common across platforms.
// It includes json tags to deserialize configuration from a file
// using the same names that the flags in the command line use.
type CommonConfig struct {
	Pidfile string `json:"pidfile,omitempty"`

	TrustKeyPath string `json:"-"`
	CorsHeaders  string `json:"api-cors-header,omitempty"`
	EnableCors   bool   `json:"api-enable-cors,omitempty"`

	Debug     bool     `json:"debug,omitempty"`
	Hosts     []string `json:"hosts,omitempty"`
	LogLevel  string   `json:"log-level,omitempty"`
	TLS       bool     `json:"tls,omitempty"`
	TLSVerify bool     `json:"tlsverify,omitempty"`

	// Embedded structs that allow config
	// deserialization without the full struct.
	CommonTLSOptions

	reloadLock sync.Mutex

	valuesSet map[string]interface{}
}

// InstallCommonFlags adds command-line options to the top-level flag parser for
// the current process.
// Subsequent calls to `flag.Parse` will populate config with values parsed
// from the command-line.
func (config *Config) InstallCommonFlags(cmd *flag.FlagSet) {

	cmd.StringVar(&config.Pidfile, "pidfile", defaultPidFile, "Path to use for daemon PID file")
	cmd.StringVar(&config.CorsHeaders, "api-cors-header", "", "Set CORS headers in the remote API")

}

// IsValueSet returns true if a configuration value
// was explicitly set in the configuration file.
func (config *Config) IsValueSet(name string) bool {
	if config.valuesSet == nil {
		return false
	}
	_, ok := config.valuesSet[name]
	return ok
}

// ReloadConfiguration reads the configuration in the host and reloads the daemon and server.
func ReloadConfiguration(configFile string, flags *flag.FlagSet, reload func(*Config)) error {
	logrus.Infof("Got signal to reload configuration, reloading from: %s", configFile)

	reload(newConfig)
	return nil
}
