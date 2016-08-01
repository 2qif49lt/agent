package daemon

import (
	"os"

	flag "github.com/2qif49lt/pflag"
)

var (
	defaultPidFile = os.Getenv("programdata") + string(os.PathSeparator) + "agentd.pid"
)

// Config defines the configuration of a docker daemon.
// These are the configuration settings that you pass
// to the docker daemon when you launch it with say: `docker daemon -e windows`
type Config struct {
	CommonConfig

	// Fields below here are platform specific. (There are none presently
	// for the Windows daemon.)
}

// InstallFlags adds command-line options to the top-level flag parser for
// the current process.
// Subsequent calls to `flag.Parse` will populate config with values parsed
// from the command-line.
func (config *Config) InstallFlags(cmd *flag.FlagSet) {
	// First handle install flags which are consistent cross-platform
	config.InstallCommonFlags(cmd)

}
