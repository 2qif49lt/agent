package cli

import (
	"github.com/2qif49lt/agent/cfg"
	flag "github.com/2qif49lt/pflag"
)

// CliFlags represents flags for the agent client.
type CliFlags struct {
	FlagSet   *flag.FlagSet
	Common    *CommonFlags
	PostParse func()

	ConfigDir string
}
