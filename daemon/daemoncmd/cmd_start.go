package daemoncmd

import (
	"github.com/2qif49lt/agent/cfg"
	"github.com/2qif49lt/agent/cli"
	"github.com/2qif49lt/cobra"
	"github.com/kardianos/service"
)

func newStartCommand() *cobra.Command {
	daemonCli := NewDaemonCli()

	cmd := &cobra.Command{
		Use:   "start [OPTIONS]",
		Short: "启动agentd",
		Args:  cli.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runStart(daemonCli)
		},
	}

	return cmd
}

func runStart(daemonCli *DaemonCli) error {
	prg := &program{daemonCli}

	if isconsole := service.Interactive(); isconsole {
		prg.StartConsole()
	} else {
		svcConfig := &service.Config{
			Name: cfg.C.SrvName,
		}
		srv, err := service.New(prg, svcConfig)
		if err != nil {
			return err
		}
		return srv.Run()
	}

	return nil
}
