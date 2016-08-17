package main

import (
	_ "github.com/2qif49lt/dump"

	"github.com/2qif49lt/agent/cfg"
	"github.com/2qif49lt/agent/cli/cmds"
)

func main() {
	comflag := cfg.InitCommonFlags()
	cobraAdaptor := cmds.NewCobraAdaptor(comflag)
	cobraAdaptor.Cmd().Execute()
}
