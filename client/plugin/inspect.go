package plugin

import (
	"encoding/json"
	"fmt"

	"github.com/2qif49lt/agent/cli"
	"github.com/2qif49lt/agent/client"
	"github.com/2qif49lt/cobra"
)

func newInspectCommand(agentCli *client.AgentCli) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "inspect",
		Short: "Inspect a plugin",
		Args:  cli.ExactArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return agentCli.Initialize()
		},

		RunE: func(cmd *cobra.Command, args []string) error {
			return runInspect(agentCli, args[0])
		},
	}

	return cmd
}

func runInspect(agentCli *client.AgentCli, name string) error {
	p, err := agentCli.Client().PluginInspect(name)
	if err != nil {
		return err
	}

	b, err := json.MarshalIndent(p, "", "\t")
	if err != nil {
		return err
	}
	fmt.Println(string(b))
	return err
}
