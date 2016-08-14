package client

import (
	"fmt"

	"golang.org/x/net/context"

	"github.com/2qif49lt/agent/cli"
	"github.com/2qif49lt/cobra"
)

// NewInfoCommand returns a cobra command for `agent info` subcommands
func NewInfoCommand(agentCli *AgentCli) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "info",
		Short: "Displays system-wide information.",
		Args:  cli.NoArgs,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return agentCli.Initialize()
		},
		Run: func(cmd *cobra.Command, args []string) {
			err := runInfo(agentCli)
			if err != nil {
				fmt.Println(err)
			}
		},
	}

	return cmd
}
func runInfo(agentCli *AgentCli) error {
	info, err := agentCli.client.Info(context.Background())
	if err != nil {
		return err
	}

	//显示结果

	fmt.Println(info)
	return err
}
