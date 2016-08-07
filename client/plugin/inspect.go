package plugin

import (
	"encoding/json"
	"fmt"

	"github.com/2qif49lt/agent/cli"
	"github.com/2qif49lt/agent/client"
	"github.com/2qif49lt/cobra"
	"golang.org/x/net/context"
)

func newInspectCommand(agentCli *client.AgentCli) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "inspect",
		Short: "Inspect a plugin",
		Args:  cli.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runInspect(agentCli, args[0])
		},
	}

	return cmd
}

func runInspect(agentCli *client.DockerCli, name string) error {
	p, err := agentCli.Client().PluginInspect(context.Background(), name)
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
