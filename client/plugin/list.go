package plugin

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/2qif49lt/agent/cli"
	"github.com/2qif49lt/agent/client"
	"github.com/2qif49lt/cobra"
	"golang.org/x/net/context"
)

// agent plugin ls
func newListCommand(agentCli *client.AgentCli) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "ls",
		Short:   "List plugins",
		Aliases: []string{"list"},
		Args:    cli.ExactArgs(0),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return agentCli.Initialize()
		},

		RunE: func(cmd *cobra.Command, args []string) error {
			return runList(agentCli)
		},
	}

	return cmd
}

func runList(agentCli *client.AgentCli) error {
	plugins, err := agentCli.Client().PluginList(context.Background())
	if err != nil {
		return err
	}

	w := tabwriter.NewWriter(os.Stdout, 20, 1, 3, ' ', 0)
	fmt.Fprintf(w, "NAME \tTAG \tACTIVE")
	fmt.Fprintf(w, "\n")

	for _, p := range plugins {
		fmt.Fprintf(w, "%s\t%s\t%v\n", p.Name, p.Tag, p.Active)
	}
	w.Flush()
	return nil
}
