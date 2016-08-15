package daemon

import (
	"fmt"
	"github.com/2qif49lt/cobra"
	"github.com/2qif49lt/logrus"
	"github.com/kardianos/service"
)

var (
	daemonFlagSrvName      = "name"
	daemonFlagSrvNameShort = "n"
)

// NewDaemonCommand creats a new cobra.Command for `agent daemon`
func NewDaemonCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "daemon",
		Short: "Manage agent daemon",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(cmd.UsageString())
		},
	}

	cmd.AddCommand(
		newInstallCommand(),
		newStartCommand(),
		newUnInstallCommand(),
	)
	cmd.PersistentFlags().StringP(daemonFlagSrvName, daemonFlagSrvNameShort, "",
		"指定服务名,若空则使用配置文件内值,若无配置则默认agentd")
	return cmd
}

type program struct {
	daemonCli *DaemonCli
}

func (p *program) Start(s service.Service) error {
	logrus.SetDefaultFileOut()

	err := p.daemonCli.start()
	if err != nil {
		return err
	}

	go p.daemonCli.run()
	return nil
}

func (p *program) StartConsole() {
	err := p.daemonCli.start()
	if err != nil {
		return
	}

	p.daemonCli.run()
}

func (p *program) Stop(s service.Service) error {
	return nil
}
