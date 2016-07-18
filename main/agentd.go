package main

import (
	"github.com/kardianos/service"
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
		"version":   Version,
		"buildtime": Buildtime,
	}).Info("Agent start!")

	for {
		if log.IsTerminal() {
			log.Println("Is Terminal")

		} else {
			l.Println("Not Terminal")
		}
		time.Sleep(time.Second / 3)

	}
}
