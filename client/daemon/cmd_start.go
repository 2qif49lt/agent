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

	cmd := &cobra.Command{
		Use:   "start [OPTIONS]",
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
			return runStart(daemonCli)
		},
	}
	flags := cmd.Flags()
	daemonCli.InitFlags(flags)
	return cmd
}

func runStart(daemonCli *DaemonCli) error {
	prg := &program{daemonCli}

	if isconsole := service.Interactive(); isconsole {
		return prg.StartConsole()
	} else {
		svcConfig := &service.Config{
			Name: daemonCli.Config.SrvName,
		}
		srv, err := service.New(prg, svcConfig)
		err = srv.Run()
		logrus.WithFields(logrus.Fields{
			"name":   svcConfig.Name,
			"return": utils.ErrStr(err),
		}).Info(`start service`)
		return err
	}
}
