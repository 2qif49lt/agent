package daemoncmd

import (
	"fmt"
	"github.com/2qif49lt/agent/daemon"
)

func newStartCommand() *cobra.Command {
	daemoncfg := &daemon.Config{}

	cmd := &cobra.Command{
		Use:   "start [OPTIONS]",
		Short: "启动agentd",
		Args:  cli.NoArgs(),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runStart(daemoncfg)
		},
	}

	flags := cmd.Flags()
	daemoncfg.InstallFlags(flags)

	return cmd
}

func runStart(daemoncfg *daemon.Config) error {
	return nil
}

func newStopCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stop",
		Short: "停止agentd",
		Args:  cli.NoArgs(),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runStop()
		},
	}

	return cmd
}

func runStop() error {
	return nil
}

func newReStartCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "restart [OPTIONS]",
		Short: "重启",
		Args:  cli.NoArgs(),
		RunE: func(cmd *cobra.Command, args []string) error {
			err := runStatus()
			if err != nil {
				return err
			}

			err = runStop()
			if err != nil {
				return err
			}

			return runStart()
		},
	}

	return cmd
}

func newStatusCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "status",
		Short: "查看agentd 运行状态",
		Args:  cli.NoArgs(),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runStatus()
		},
	}

	return cmd
}

func runStatus() error {
	return nil
}
