package daemon

import (
	"github.com/2qif49lt/agent/cfg"
	"github.com/2qif49lt/agent/cli"
	coredaemon "github.com/2qif49lt/agent/daemon"
	"github.com/2qif49lt/agent/utils"
	"github.com/2qif49lt/cobra"
	"github.com/2qif49lt/logrus"
	"github.com/kardianos/service"
)

func newInstallCommand() *cobra.Command {
	conf := &coredaemon.Config{}

	install := &cobra.Command{
		Use:   "install",
		Short: "本地功能,将安装agentd服务",
		Args:  cli.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			name := cfg.Conf.SrvName
			if conf.SrvName != "" {
				name = conf.SrvName
			}
			//		conargs := make([]string, 0)

			logrus.Println("args", cmd.Name(), cmd.PersistentFlags().Args())
			return runInstall(name, args)
		},
	}
	flags := install.Flags()
	conf.InstallFlags(flags)
	return install
}

func runInstall(name string, args []string) error {
	srvargs := append([]string{"daemon", "start"}, args...)

	svcConfig := &service.Config{
		Name:        name,
		DisplayName: name,
		Description: `Night gathers, and now my watch begins. It shall not end until my death.`,
		Arguments:   srvargs,
		Option:      service.KeyValue{"RunAtLoad": true},
	}
	prg := &program{}
	srv, err := service.New(prg, svcConfig)
	if err != nil {
		return err
	}
	err = srv.Install()
	logrus.WithFields(logrus.Fields{
		"name":   name,
		"return": utils.ErrStr(err),
	}).Info(`install service`)
	return err
}

func newUnInstallCommand() *cobra.Command {
	sn := ""

	cmd := &cobra.Command{
		Use:   "uninstall",
		Short: "本地功能,将卸载agentd服务",
		Args:  cli.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			name := cfg.Conf.SrvName
			snflag := cmd.Flags().Lookup("name")
			if snflag != nil {
				sn := snflag.Value.String()
				if sn != "" {
					name = sn
				}
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
	err = srv.Uninstall()
	logrus.WithFields(logrus.Fields{
		"name":   name,
		"return": utils.ErrStr(err),
	}).Info(`uninstall service`)
	return err
}
