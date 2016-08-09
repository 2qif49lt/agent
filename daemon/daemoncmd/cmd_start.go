package daemoncmd

import (
	"fmt"
	"github.com/2qif49lt/agent/cfg"
	"github.com/2qif49lt/agent/cli"
	"github.com/2qif49lt/agent/daemon"
	"github.com/2qif49lt/cobra"
)

const (
	daemonConfigFileFlag = "config-file"
)

func newStartCommand() *cobra.Command {
	daemonCli := NewDaemonCli()

	cmd := &cobra.Command{
		Use:   "start [OPTIONS]",
		Short: "启动agentd",
		Args:  cli.NoArgs(),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runStart(daemonCli)
		},
	}

	flags := cmd.Flags()
	daemonCli.InstallFlags(flags)
	defaultConfigPath := filepath.Join(cfg.GetConfigDir(), cfg.ConfigFileName)

	daemonCli.configFile = flags.String(daemonConfigFileFlag, defaultConfigPath,
		fmt.Sprintf("Daemon configuration file,default: %s"), defaultConfigPath)

	return cmd
}

func runStart(daemonCli *DaemonCli) error {
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
