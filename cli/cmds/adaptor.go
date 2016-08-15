package cmds

import (
	"github.com/2qif49lt/agent/cfg"
	"github.com/2qif49lt/agent/client"
	"github.com/2qif49lt/agent/client/daemon"
	"github.com/2qif49lt/agent/client/plugin"
	"github.com/2qif49lt/agent/client/system"

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
		Use:   "agent",
		Short: "A self-sufficient DevOps Agent.",

		SilenceUsage:  true,
		SilenceErrors: false,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return cfg.PostCheck()
		},
	}

	rootCmd.SetOutput(os.Stdout)
	rootCmd.AddCommand(
		system.NewEventsCommand(agentCli),
		system.NewVersionCommand(agentCli),
		daemon.NewDaemonCommand(),
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
