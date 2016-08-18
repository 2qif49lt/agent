package daemon

import (
	"github.com/2qif49lt/agent/cfg"
	"github.com/2qif49lt/agent/cli"
	"github.com/2qif49lt/agent/utils"
	"github.com/2qif49lt/cobra"
	"github.com/2qif49lt/logrus"
	"github.com/kardianos/service"
)

func newStartCommand() *cobra.Command {
	daemonCli := NewDaemonCli()
	console := false

	cmd := &cobra.Command{
		Use:   "start",
		Short: "本地功能,启动agentd",
		Args:  cli.NoArgs,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if daemonCli.Config.SrvName == "" {
				daemonCli.Config.SrvName = cfg.Conf.SrvName
			}
			if daemonCli.Config.AgentID == "" {
				daemonCli.Config.AgentID = cfg.Conf.Agentid
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return runStart(daemonCli, console)
		},
	}
	flags := cmd.Flags()
	daemonCli.InitFlags(flags)
	flags.BoolVarP(&console, "console", "c", false, "是否cosole模式")

	return cmd
}

func runStart(daemonCli *DaemonCli, console bool) error {
	prg := &program{daemonCli}

	if console {
		return prg.StartConsole()
	} else {
		svcConfig := &service.Config{
			Name: daemonCli.Config.SrvName,
		}
		srv, _ := service.New(prg, svcConfig)

		var err error
		if service.Interactive() == false {
			err = srv.Run()
		} else {
			err = service.Control(srv, "start")
		}

		logrus.WithFields(logrus.Fields{
			"name":  svcConfig.Name,
			"error": utils.ErrStr(err),
		}).Info(`start service`)

		return err
	}
}

func newStopCommand() *cobra.Command {
	sn := ""

	cmd := &cobra.Command{
		Use:   "stop",
		Short: "本地功能,停止agentd,可以附加原因",
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

			return runStop(name)
		},
	}
	flags := cmd.Flags()
	flags.StringVarP(&sn, "name", "n", "", "指定服务名,若空则使用配置文件内值")
	return cmd
}

func runStop(name string) error {
	svcConfig := &service.Config{
		Name: name,
	}
	prg := &program{}
	srv, err := service.New(prg, svcConfig)
	if err != nil {
		return err
	}
	err = srv.Stop()
	logrus.WithFields(logrus.Fields{
		"name":  name,
		"error": utils.ErrStr(err),
	}).Info(`stop service`)
	return err
}

func newReStartCommand() *cobra.Command {
	sn := ""

	cmd := &cobra.Command{
		Use:   "restart [reason]",
		Short: "本地功能,重启agentd",
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

			return runReStart(name)
		},
	}
	flags := cmd.Flags()
	flags.StringVarP(&sn, "name", "n", "", "指定服务名,若空则使用配置文件内值")
	return cmd
}

func runReStart(name string) error {
	svcConfig := &service.Config{
		Name: name,
	}
	prg := &program{}
	srv, err := service.New(prg, svcConfig)
	if err != nil {
		return err
	}
	err = srv.Restart()
	logrus.WithFields(logrus.Fields{
		"name":  name,
		"error": utils.ErrStr(err),
	}).Info(`restart service`)
	return err
}
