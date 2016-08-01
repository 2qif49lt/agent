// +build linux freebsd

package daemon

import (
	"fmt"
	"net"

	flag "github.com/2qif49lt/pflag"
)

var (
	defaultPidFile = "/var/run/agentd.pid"
)

// Config defines the configuration of a docker daemon.
// It includes json tags to deserialize configuration from a file
// using the same names that the flags in the command line uses.
type Config struct {
	CommonConfig
}

// InstallFlags adds command-line options to the top-level flag parser for
// the current process.
// Subsequent calls to `flag.Parse` will populate config with values parsed
// from the command-line.
func (config *Config) InstallFlags(cmd *flag.FlagSet) {
	// First handle install flags which are consistent cross-platform
	config.InstallCommonFlags(cmd)
}
