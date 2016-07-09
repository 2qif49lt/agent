package main

import (
	log "github.com/Sirupsen/logrus"
)

func main() {
	log.WithFields(log.Fields{
		"version":   Version,
		"buildtime": Buildtime,
	}).Info("Agent start!")
}
