package main

import (
	"github.com/2qif49lt/agent/cfg"
	"github.com/2qif49lt/agent/version"
	log "github.com/2qif49lt/logrus"
	flag "github.com/2qif49lt/pflag"
	"github.com/kardianos/service"

	"time"
)

type program struct{}

func (p *program) Start(s service.Service) error {
	// Start should not block. Do the actual work async.
	go p.run()
	return nil
}
func (p *program) run() {
	// Do work here
	Agentd()
}
func (p *program) Stop(s service.Service) error {
	// Stop should not block. Return with a few seconds.
	// 用于daemon退出清理
	return nil
}

func Agentd() {
	l := log.NewSSLog("log", "base.txt", log.InfoLevel)
	l.WithFields(log.Fields{
		"version":   version.SRV_VERSION,
		"buildtime": version.BUILDTIME,
	}).Info("Agentd start!")

	daemonCli = NewDaemonCli()

	flHelp = flag.BoolP("help", "h", false, "Print usage")
	flVersion = flag.BoolP("version", "v", false, "Print version information and quit")

	com := cfg.InitCommonFlags()

	flag.Merge(flag.CommandLine, com.FlagSet)

	flag.Usage = func() {
		fmt.Fprint(stdout, "Usage: agent [ --help | -h | -v | --version ]\n\n")

		flag.CommandLine.SetOutput(stdout)
		flag.PrintDefaults()
	}

	if err := flag.CommandLine.Parse(os.Args[1:]); err != nil {
		os.Exit(1)
	}

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

}
