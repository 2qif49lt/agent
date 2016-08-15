package main

import (
	_ "github.com/2qif49lt/dump"

	"github.com/2qif49lt/agent/cfg"
	"github.com/2qif49lt/agent/cli/cmds"
	"github.com/2qif49lt/logrus"
)

func main() {
	comflag := cfg.InitCommonFlags()
	_, err := cfg.Load()
	if err != nil {
		logrus.Panicf(`config: %s`, err.Error())
	}
	cobraAdaptor := cmds.NewCobraAdaptor(comflag)
	cobraAdaptor.Cmd().Execute()
}
