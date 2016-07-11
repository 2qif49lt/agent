package main

import (
	"github.com/2qif49lt/agent/pkg/signal"
	log "github.com/2qif49lt/logrus"
	"time"
)

func main() {
	defer func() {
		if err := recover(); err != nil {
			log.Warnln("damn", err)
			signal.DumpStacks()
		}
	}()

	signal.Trap(func() {
		log.Warnln("oops")
	})

	log.WithFields(log.Fields{
		"version":   Version,
		"buildtime": Buildtime,
	}).Info("Agent start!")

	time.Sleep(time.Second * 1)

}
