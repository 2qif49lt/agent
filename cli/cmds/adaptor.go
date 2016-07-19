package cmds

import (
	"github.com/spf13/cobra"
)

// CobraAdaptor is an adaptor for supporting spf13/cobra commands
type CobraAdaptor struct {
	rootCmd   *cobra.Command
	dockerCli *client.DockerCli
}
