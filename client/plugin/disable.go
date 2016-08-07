package plugin

import (
	"github.com/2qif49lt/agent/cli"
	"github.com/2qif49lt/agent/client"
	"github.com/2qif49lt/cobra"
	"golang.org/x/net/context"
)

// agent plugin disable PLUGIN

func newDisableCommand(agentCli *client.AgentCli) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "disable",
		Short: "Disable a plugin",
		Args:  cli.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return agentCli.Client().PluginDisable(context.Background(), args[0])
		},
	}

	return cmd
}
