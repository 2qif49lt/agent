package plugin

import (
	"github.com/2qif49lt/agent/cli"
	"github.com/2qif49lt/agent/client"
	"github.com/2qif49lt/cobra"
	//	"context"
	//从软件中心下载
	//	"github.com/docker/docker/reference"
	//	"github.com/docker/docker/registry"
)

// agent plugin install name/version --disable
func newInstallCommand(agentCli *client.AgentCli) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "install",
		Short: "Install a plugin",
		Args:  cli.RequiresMinArgs(1), // 插件名,参数为版本号
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return agentCli.Initialize()
		},

		RunE: func(cmd *cobra.Command, args []string) error {
			return runInstall(agentCli, args[0], args[1:])
		},
	}

	return cmd
}

func runInstall(agentCli *client.AgentCli, name string, args []string) error {

	//	ctx := context.Background()

	// TODO: pass acceptAllPermissions and noEnable flag
	//	return agentCli.Client().PluginInstall(ctx, ref.String(), encodedAuth)

	return nil
}
