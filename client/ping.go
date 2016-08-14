package client

import (
	"fmt"
	"strings"

	"golang.org/x/net/context"

	"github.com/2qif49lt/agent/cli"
	"github.com/2qif49lt/cobra"
)

// NewPingCommand returns a cobra command for `agent ping some message` subcommands
func NewPingCommand(agentCli *AgentCli) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "ping [message ...]",
		Short:   "Ping agentd,receive a pong.",
		Args:    cli.RequiresMinArgs(1),
		Example: "* ping hello",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return agentCli.Initialize()
		},
		Run: func(cmd *cobra.Command, args []string) {
			err := runPing(agentCli, args)
			if err != nil {
				fmt.Println(err)
			}
		},
	}
	return cmd
}
func runPing(agentCli *AgentCli, args []string) error {
	pong, err := agentCli.client.Ping(context.Background(), strings.Join(args, " "))
	if err != nil {
		return err
	}

	//显示结果

	fmt.Println(args, pong)
	return err
}
