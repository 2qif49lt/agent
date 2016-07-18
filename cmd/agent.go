package main

import (
	_ "github.com/2qif49lt/dump"
	log "github.com/2qif49lt/logrus"

	"github.com/2qif49lt/agent/pkg/signal"
	"github.com/2qif49lt/cfg"
	flag "github.com/2qif49lt/pflag"
	"time"
)

var (
	comflag   = cfg.InitCommonFlags()
	flHelp    = flag.BoolP("help", "h", false, "Print usage")
	flVersion = flag.BoolP("version", "v", false, "Print version information and quit")
)

func main() {

	flag.Merge(flag.CommandLine, cfg.commonFlags.FlagSet)
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
	if utils.ExperimentalBuild() {
		fmt.Printf("Docker version %s, build %s, experimental\n", dockerversion.Version, dockerversion.GitCommit)
	} else {
		fmt.Printf("Docker version %s, build %s\n", dockerversion.Version, dockerversion.GitCommit)
	}
}
