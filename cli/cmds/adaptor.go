package cmds

import (
	"github.com/2qif49lt/agent/api/client"
	"github.com/2qif49lt/cobra"
)

// CobraAdaptor is an adaptor for supporting spf13/cobra commands
type CobraAdaptor struct {
	rootCmd   *cobra.Command
	dockerCli *client.DockerCli
}

// NewCobraAdaptor returns a new handler
func NewCobraAdaptor(clientFlags *cliflags.ClientFlags) CobraAdaptor {
	stdin, stdout, stderr := term.StdStreams()
	dockerCli := client.NewDockerCli(stdin, stdout, stderr, clientFlags)

	var rootCmd = &cobra.Command{
		Use:           "agent",
		SilenceUsage:  true,
		SilenceErrors: true,
	}
	rootCmd.SetUsageTemplate(usageTemplate)
	rootCmd.SetHelpTemplate(helpTemplate)
	rootCmd.SetFlagErrorFunc(cli.FlagErrorFunc)
	rootCmd.SetOutput(stdout)
	rootCmd.AddCommand(
		system.NewEventsCommand(dockerCli),
		system.NewVersionCommand(dockerCli),
	)
	plugin.NewPluginCommand(rootCmd, dockerCli)

	rootCmd.PersistentFlags().BoolP("help", "h", false, "Print usage")
	rootCmd.PersistentFlags().MarkShorthandDeprecated("help", "please use --help")

	return CobraAdaptor{
		rootCmd:   rootCmd,
		dockerCli: dockerCli,
	}
}
