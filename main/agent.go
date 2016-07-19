package main

import (
	_ "github.com/2qif49lt/dump"
	log "github.com/2qif49lt/logrus"

	"github.com/2qif49lt/agent/pkg/signal"
	"github.com/2qif49lt/cfg"
	flag "github.com/2qif49lt/pflag"
	"os"
	"time"
)

var (
	comflag   = cfg.InitCommonFlags()
	flHelp    = flag.BoolP("help", "h", false, "Print usage")
	flVersion = flag.BoolP("version", "v", false, "Print version information and quit")
)

func main() {

	flag.Merge(flag.CommandLine, comflag.FlagSet)

	flag.Usage = func() {
		fmt.Fprint(os.Stdout, "Usage: agent [OPTIONS] COMMAND [arg...]\n       agent [ --help |-h | -v | --version ]\n\n")
		fmt.Fprint(os.Stdout, "A self-sufficient runtime tool.\n\nOptions:\n")

		flag.CommandLine.SetOutput(stdout)
		flag.PrintDefaults()

		help := "\nCommands:\n"

		dockerCommands := append(cli.DockerCommandUsage, cobraAdaptor.Usage()...)
		for _, cmd := range sortCommands(dockerCommands) {
			help += fmt.Sprintf("    %-10.10s%s\n", cmd.Name, cmd.Description)
		}

		help += "\nRun 'docker COMMAND --help' for more information on a command."
		fmt.Fprintf(stdout, "%s\n", help)
	}

	flag.Parse()

	if *flVersion {
		showVersion()
		return
	}

	if *flHelp {
		// if global flag --help is present, regardless of what other options and commands there are,
		// just print the usage.
		flag.Usage()
		return
	}

	if flag.IsSet("daemon") {

	}
	svcConfig := &service.Config{
		Name:        "GoServiceExampleSimple",
		DisplayName: "Go Service Example",
		Description: "This is an example Go service.",
	}

	prg := &program{}
	s, err := service.New(prg, svcConfig)
	if err != nil {
		log.Fatal(err)
	}

	err = s.Run()
	if err != nil {
		log.Error(err)
	}

}

func Agentc() {
	signal.Trap(func() {
		println("trap")
	})

}

func showVersion() {
	fmt.Printf("Agent version %s, build %s\n", VERSION, BUILDTIME)
}
