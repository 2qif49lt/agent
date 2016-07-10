package main

import (
	log "github.com/2qif49lt/logrus"
	"time"
)

func main() {
	for i := 0; i < 1000; i++ {
		log.WithFields(log.Fields{
			"version":   Version,
			"buildtime": Buildtime,
		}).Info("Agent start!")
		time.Sleep(time.Second / 3)
	}
}
