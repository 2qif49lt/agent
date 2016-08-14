package cmds

import (
	"github.com/2qif49lt/agent/cfg"
	"github.com/2qif49lt/agent/client"
	"github.com/2qif49lt/agent/client/plugin"
	"github.com/2qif49lt/agent/client/system"
	"github.com/2qif49lt/agent/daemon/daemoncmd"

	"github.com/2qif49lt/cobra"
	"github.com/2qif49lt/pflag"

	"os"
)

// CobraAdaptor is an adaptor for supporting spf13/cobra commands
type CobraAdaptor struct {
	rootCmd  *cobra.Command
	agentCli *client.AgentCli
}

func (adaptor *CobraAdaptor) Cmd() *cobra.Command {
	return adaptor.rootCmd
}

// NewCobraAdaptor returns a new handler
func NewCobraAdaptor(com *cfg.CommonFlags) CobraAdaptor {
	agentCli := client.NewAgentCli(com)

	var rootCmd = &cobra.Command{
		Use:           "agent",
		SilenceUsage:  true,
		SilenceErrors: false,
	}

	rootCmd.SetOutput(os.Stdout)
	rootCmd.AddCommand(
		system.NewEventsCommand(agentCli),
		system.NewVersionCommand(agentCli),
		daemoncmd.NewDaemonCommand(),
		client.NewInfoCommand(agentCli),
		client.NewPingCommand(agentCli),
	)
	plugin.NewPluginCommand(rootCmd, agentCli)

	cmdflags := rootCmd.PersistentFlags()
	pflag.Merge(cmdflags, com.FlagSet)

	return CobraAdaptor{
		rootCmd:  rootCmd,
		agentCli: agentCli,
	}
}
