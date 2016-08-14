package plugin

import (
	"github.com/2qif49lt/agent/cli"
	"github.com/2qif49lt/agent/client"
	"github.com/2qif49lt/cobra"
	"golang.org/x/net/context"
)

func newEnableCommand(agentCli *client.AgentCli) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "enable",
		Short: "Enable a plugin",
		Args:  cli.ExactArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return agentCli.Initialize()
		},

		RunE: func(cmd *cobra.Command, args []string) error {
			return agentCli.Client().PluginEnable(context.Background(), args[0])
		},
	}

	return cmd
}
