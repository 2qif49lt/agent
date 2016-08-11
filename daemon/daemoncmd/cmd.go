package daemoncmd

import (
	"fmt"
	"github.com/2qif49lt/cobra"
	"github.com/2qif49lt/logrus"
	"github.com/kardianos/service"
)

// NewDaemonCommand creats a new cobra.Command for `agent daemon`
func NewDaemonCommand() *cobra.Command {
	install := false
	uninstall := false
	interact := false

	cmd := &cobra.Command{
		Use:   "daemon ",
		Short: "Manage agent daemon",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(cmd.UsageString())
		},
	}

	flags := cmd.Flags()

	flags.BoolVar(&install, "install", false, "安装 agentd 作为服务")
	flags.BoolVar(&uninstall, "uninstall", false, "进行卸载 agentd")
	flags.BoolVar(&interact, "interact", false, "控制台模式")

	cmd.AddCommand(
		newInstallCommand(),
		newStartCommand(),
		newUnInstallCommand(),
	)
	return cmd
}

/*
	简化成 agent daemon -i / -d / -u / -s
*/

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
