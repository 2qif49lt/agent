package client

import (
	"fmt"
	"net/url"
	"strings"

	"golang.org/x/net/context"

	Cli "github.com/2qif49lt/agent/cli"
	flag "github.com/2qif49lt/pflag"
)

// CmdPing test agentd,receive a pong
//
// Usage: agent info
func (cli *AgentCli) CmdPing(args ...string) error {
	cmd := Cli.Subcmd("ping")
	cmd.Require(flag.Min, 1)

	cmd.ParseFlags(args)

	pong, err := cli.client.Ping(context.Background(), strings.Join(args, " "))
	if err != nil {
		return err
	}

	//显示结果

	fmt.Println(pong)
	return nil
}
