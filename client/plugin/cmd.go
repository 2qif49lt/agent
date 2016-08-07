package plugin

import (
	"fmt"

	"github.com/2qif49lt/agent/cli"
	"github.com/2qif49lt/agent/client"
	"github.com/2qif49lt/cobra"
)

// NewPluginCommand returns a cobra command for `plugin` subcommands
func NewPluginCommand(rootCmd *cobra.Command, agentCli *client.AgentCli) {
	cmd := &cobra.Command{
		Use:   "plugin",
		Short: "Manage agent plugins",
		Args:  cli.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(cmd.UsageString())
		},
	}

	cmd.AddCommand(
		newDisableCommand(agentCli),
		newEnableCommand(agentCli),
		newInspectCommand(agentCli),
		newInstallCommand(agentCli),
		newListCommand(agentCli),
		newRemoveCommand(agentCli),
	)

	rootCmd.AddCommand(cmd)
}
