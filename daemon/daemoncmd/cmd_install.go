package daemoncmd

import (
	"fmt"
	"github.com/2qif49lt/agent/cfg"
	"github.com/2qif49lt/agent/cli"
	"github.com/2qif49lt/cobra"
	"github.com/kardianos/service"
)

func newInstallCommand() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "install [service name]",
		Short: "安装 agentd 作为服务,若不提供名称则以配置文件内 srvname",
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
	svcConfig := &service.Config{
		Name:        name,
		DisplayName: name,
		Description: fmt.Sprintf(`i am %s`, name),
		Arguments:   []string{"daemon", "start"},
	}
	prg := &program{}
	srv, err := service.New(prg, svcConfig)
	if err != nil {
		return err
	}
	return srv.Install()
}

func newUnInstallCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "uninstall [service name]",
		Short: "卸载 agentd 服务,若不提供名称则以配置文件内 srvname",
		Args:  cli.RequiresMaxArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := cfg.C.SrvName
			if len(args) > 0 {
				name = args[0]
			}
			return runUnInstall(name)
		},
	}

	return cmd
}

func runUnInstall(name string) error {
	svcConfig := &service.Config{
		Name: name,
	}
	prg := &program{}
	srv, err := service.New(prg, svcConfig)
	if err != nil {
		return err
	}
	return srv.Uninstall()

	return nil
}
