package plugin

import (
	"fmt"

	"github.com/2qif49lt/agent/cli"
	"github.com/2qif49lt/agent/client"
	"github.com/2qif49lt/cobra"
	"golang.org/x/net/context"
)

// agent plugin rm PLUGIN
func newRemoveCommand(agentCli *client.AgentCli) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "rm",
		Short:   "Remove a plugin",
		Aliases: []string{"remove"},
		Args:    cli.RequiresMinArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return agentCli.Initialize()
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return runRemove(agentCli, args)
		},
	}

	return cmd
}

func runRemove(agentCli *client.AgentCli, names []string) error {
	for _, name := range names {
		fmt.Printf("%s:", name)

		// TODO: pass names to api instead of making multiple api calls
		err := agentCli.Client().PluginRemove(context.Background(), name)
		if err == nil {
			fmt.Printf("OK")
		} else {
			fmt.Printf("FAIL")
			return err
		}

	}
	return nil
}
