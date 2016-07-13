package main

import (
	log "github.com/2qif49lt/logrus"

	"github.com/2qif49lt/agent/pkg/signal"
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
	return nil
}

func main() {

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

func Agentd() {
	l := log.NewSSLog("log", "base.txt", log.InfoLevel)
	signal.Trap(func() {
		l.Warnln("exit by signal")
	})

	log.WithFields(log.Fields{
		"version":   Version,
		"buildtime": Buildtime,
	}).Info("Agent start!")

	for i := 0; i < 100; i++ {
		if log.IsTerminal() {
			log.Println("Is Terminal")

		} else {
			log.Println("Not Terminal")
		}
		time.Sleep(time.Second / 10)

	}
	i := 0
	j := 10 / i
	_ = j
}
