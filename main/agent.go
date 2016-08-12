package main

import (
	_ "github.com/2qif49lt/dump"
	//	log "github.com/2qif49lt/logrus"

	"github.com/2qif49lt/agent/cfg"
	"github.com/2qif49lt/agent/cli"
	"github.com/2qif49lt/agent/cli/cmds"
	"github.com/2qif49lt/agent/client"
	"github.com/2qif49lt/agent/version"
	//"github.com/2qif49lt/agent/pkg/signal"
	flag "github.com/2qif49lt/pflag"

	"fmt"
	"os"
	//	"time"
)

var (
	comflag = cfg.InitCommonFlags()
	//	flHelp    = flag.BoolP("help", "h", false, "Print usage")
	flVersion = flag.BoolP("version", "v", false, "Print version information and quit")
)

func main() {

	flag.Merge(flag.CommandLine, comflag.FlagSet)
	cobraAdaptor := cmds.NewCobraAdaptor(comflag)

	flag.Usage = func() {
		fmt.Fprint(os.Stdout, "Usage: agent [OPTIONS] COMMAND [arg...]\n       agent [ --help |-h | -v | --version ]\n\n")
		fmt.Fprint(os.Stdout, "\nOptions:\n")

		flag.CommandLine.SetOutput(os.Stdout)
		flag.PrintDefaults()

		help := "\nCommands:\n"
		for _, cmd := range cobraAdaptor.Usage() {
			help += fmt.Sprintf("    %-10.10s%s\n", cmd.Name, cmd.Description)
		}

		help += "\nRun 'agent COMMAND --help' for more information on a command."
		fmt.Fprintf(os.Stdout, "%s\n", help)
	}

	flag.Parse()

	if *flVersion {
		showVersion()
		return
	}

	clientCli := client.NewAgentCli(comflag)

	c := cli.New(clientCli, cobraAdaptor)
	if err := c.Run(flag.Args()...); err != nil {
		if sterr, ok := err.(cli.StatusError); ok {
			if sterr.Status != "" {
				fmt.Fprintln(os.Stdout, sterr.Status)
			}
			// StatusError should only be used for errors, and all errors should
			// have a non-zero exit status, so never exit with 0
			if sterr.StatusCode == 0 {
				os.Exit(1)
			}
			os.Exit(sterr.StatusCode)
		}
		fmt.Fprintln(os.Stdout, err)
		os.Exit(1)
	}
}

func showVersion() {
	fmt.Printf("Agent version %s, build %s\n", version.CLI_VERSION, version.BUILDTIME)
}
