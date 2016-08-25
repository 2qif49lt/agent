package daemon

import (
	"fmt"
	"github.com/2qif49lt/cobra"
	"github.com/2qif49lt/logrus"
	"github.com/kardianos/service"
)

// NewDaemonCommand creats a new cobra.Command for `agent daemon`
func NewDaemonCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "daemon",
		Short: "Manage agent daemon.",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(cmd.UsageString())
		},
	}

	cmd.AddCommand(
		newInstallCommand(),
		newStartCommand(),
		newUnInstallCommand(),
		newStopCommand(),
		newReStartCommand(),
	)
	// opt-json为client 专属
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

func (p *program) StartConsole() error {
	err := p.daemonCli.start()
	if err != nil {
		return err
	}

	p.daemonCli.run()
	return nil
}

func (p *program) Stop(s service.Service) error {
	return nil
}
