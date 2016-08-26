package cmds

import (
	"github.com/2qif49lt/agent/cfg"
	"github.com/2qif49lt/agent/client"
	"github.com/2qif49lt/agent/client/daemon"
	"github.com/2qif49lt/agent/client/plugin"
	"github.com/2qif49lt/agent/client/system"

	"github.com/2qif49lt/cobra"
	flag "github.com/2qif49lt/pflag"

	"os"
)

// CobraAdaptor is an adaptor for supporting spf13/cobra commands
type CobraAdaptor struct {
	rootCmd  *cobra.Command
	agentCli *client.AgentCli
}

func (adaptor *CobraAdaptor) Execute() error {
	return adaptor.rootCmd.Execute()
}

func (adaptor *CobraAdaptor) SetArgs(a []string) {
	adaptor.rootCmd.SetArgs(a)
}

func (adaptor *CobraAdaptor) Cmd() *cobra.Command {
	return adaptor.rootCmd
}

var (
	helptemplate = `{{.Short}}{{.Long}}
Usage:{{if .Runnable}}
  {{if .HasAvailableFlags}}{{appendIfNotPresent .UseLine "[flags]"}}{{else}}{{.UseLine}}{{end}}{{end}}{{if .HasAvailableSubCommands}}
  {{ .CommandPath}} [command]{{end}}{{if gt .Aliases 0}}

Aliases:
  {{.NameAndAliases}}
{{end}}{{if .HasExample}}

Examples:
{{ .Example }}{{end}}{{ if .HasAvailableSubCommands}}

Available Commands:{{range .Commands}}{{if .IsAvailableCommand}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{ if .HasAvailableLocalFlags}}

Flags:
{{.LocalFlags.FlagUsages | trimRightSpace}}{{end}}{{ if .HasAvailableInheritedFlags}}

Global Flags:
{{.InheritedFlags.FlagUsages | trimRightSpace}}{{end}}{{if .HasHelpSubCommands}}

Additional help topics:{{range .Commands}}{{if .IsHelpCommand}}
  {{rpad .CommandPath .CommandPathPadding}} {{.Short}}{{end}}{{end}}{{end}}{{ if .HasAvailableSubCommands }}

Use "{{.CommandPath}} [command] --help" for more information about a command.{{end}}
{{if not .HasParent}}You also can choose another way:Use only one mission file instead of entering all flags.
Usage:
  agent [flag]
Example: 
  agent -m=./mission.json

Flags:
  -m, --mission-file string   specify the json formated mission file path,the content's field will overwrite the flags{{end}}
`
)

// NewCobraAdaptor returns a new handler
func NewCobraAdaptor(com *cfg.CommonFlags) CobraAdaptor {
	agentCli := client.NewAgentCli(com)

	var rootCmd = &cobra.Command{
		Use:           "agent",
		Short:         "A self-sufficient DevOps Agent.",
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

	persinFlags := rootCmd.PersistentFlags()
	flag.Merge(persinFlags, com.FlagSet)
	agentCli.InitFlags(persinFlags)

	rootCmd.SetHelpTemplate(helptemplate)
	return CobraAdaptor{
		rootCmd:  rootCmd,
		agentCli: agentCli,
	}
}
