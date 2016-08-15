package daemon

import (
	"fmt"
	"github.com/2qif49lt/agent/cfg"
	"github.com/2qif49lt/agent/cli"
	"github.com/2qif49lt/cobra"
	"github.com/kardianos/service"
)

func newInstallCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "install",
		Short: "本地功能,将安装agentd服务",
		Args:  cli.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			name := cfg.Conf.SrvName
			snflag := cmd.Flags().Lookup(daemonFlagSrvName)
			if snflag != nil {
				sn := snflag.Value.String()
				if sn != "" {
					name = sn
				}
			}
			if name == "" {
				name = "agentd"
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
	sn := ""

	cmd := &cobra.Command{
		Use:   "uninstall",
		Short: "本地功能,将卸载agentd服务",
		Args:  cli.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			name := cfg.Conf.SrvName
			snflag := cmd.Flags().Lookup(daemonFlagSrvName)
			if snflag != nil {
				sn := snflag.Value.String()
				if sn != "" {
					name = sn
				}
			}
			if name == "" {
				name = "agentd"
			}
			return runUnInstall(name)
		},
	}
	flags := cmd.Flags()
	flags.StringVarP(&sn, "name", "n", "", "指定服务名,若空则使用配置文件内值")
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
