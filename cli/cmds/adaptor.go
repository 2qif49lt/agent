package cmds

import (
	"github.com/2qif49lt/agent/cfg"
	"github.com/2qif49lt/agent/cli"
	"github.com/2qif49lt/agent/client"
	"github.com/2qif49lt/agent/client/plugin"
	"github.com/2qif49lt/agent/client/system"
	"github.com/2qif49lt/agent/daemon/daemoncmd"

	"github.com/2qif49lt/cobra"

	"os"
)

// CobraAdaptor is an adaptor for supporting spf13/cobra commands
type CobraAdaptor struct {
	rootCmd  *cobra.Command
	agentCli *client.AgentCli
}

// NewCobraAdaptor returns a new handler
func NewCobraAdaptor(com *cfg.CommonFlags) CobraAdaptor {
	agentCli := client.NewAgentCli(com)

	var rootCmd = &cobra.Command{
		Use:           "agent",
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	rootCmd.SetOutput(os.Stdout)
	rootCmd.AddCommand(
		system.NewEventsCommand(agentCli),
		system.NewVersionCommand(agentCli),
		daemoncmd.NewDaemonCommand(),
	)
	plugin.NewPluginCommand(rootCmd, agentCli)

	rootCmd.PersistentFlags().BoolP("help", "h", false, "Print usage")

	return CobraAdaptor{
		rootCmd:  rootCmd,
		agentCli: agentCli,
	}
}

// Usage returns the list of commands and their short usage string for
// all top level cobra commands.
func (c CobraAdaptor) Usage() []cli.Command {
	cmds := []cli.Command{}
	for _, cmd := range c.rootCmd.Commands() {
		if cmd.Name() != "" {
			cmds = append(cmds, cli.Command{Name: cmd.Name(), Description: cmd.Short})
		}
	}
	return cli.SortCommands(cmds)
}

func (c CobraAdaptor) run(cmd string, args []string) error {
	if err := c.agentCli.Initialize(); err != nil {
		return err
	}
	// Prepend the command name to support normal cobra command delegation
	c.rootCmd.SetArgs(append([]string{cmd}, args...))
	return c.rootCmd.Execute()
}

// Command returns a cli command handler if one exists
func (c CobraAdaptor) Command(name string) func(...string) error {
	for _, cmd := range c.rootCmd.Commands() {
		if cmd.Name() == name {
			return func(args ...string) error {
				return c.run(name, args)
			}
		}
	}
	return nil
}
