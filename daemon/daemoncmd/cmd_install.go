package daemoncmd

import (
	"fmt"
	"github.com/2qif49lt/agent/cfg"
	"github.com/2qif49lt/agent/cli"
	"github.com/2qif49lt/agent/daemon"
	"github.com/2qif49lt/cobra"
)

func newInstallCommand() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "install [service name]",
		Short: "安装 agentd 作为服务",
		Args:  cli.RequiresMaxArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := cfg.C.SrvName
			if len(args) > 0 {
				name = args[0]
			}
			return runInstall(name)
		},
	}

	return cmd
}

func runInstall(name string) error {
	return nil
}

func newUnInstallCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "uninstall",
		Short: "安装 agentd 作为服务",
		Args:  cli.NoArgs(),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runUnInstall(name)
		},
	}

	return cmd
}

func runUnInstall(name string) error {
	return nil
}
