package client

import (
	"fmt"

	"golang.org/x/net/context"

	Cli "github.com/2qif49lt/agent/cli"
	flag "github.com/2qif49lt/pflag"
)

// CmdInfo displays system-wide information.
//
// Usage: agent info
func (cli *AgentCli) CmdInfo(args ...string) error {
	cmd := Cli.Subcmd("info")
	cmd.Require(flag.Exact, 0)

	cmd.ParseFlags(args)

	info, err := cli.client.Info(context.Background())
	if err != nil {
		return err
	}

	//显示结果

	fmt.Println(info)
	return nil
}
