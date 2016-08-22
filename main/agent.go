package main

import (
	_ "github.com/2qif49lt/dump"

	"github.com/2qif49lt/agent/cfg"
	"github.com/2qif49lt/agent/cli/cmds"
	"github.com/2qif49lt/agent/utils"

	"math/rand"
	"os"
	"time"
)

func initProgram() {
	rand.Seed(time.Now().Unix())
	procdir, _ := utils.GetProcAbsDir()
	os.Chdir(procdir)
}

func init() {
	initProgram()
}

func main() {

	comflag := cfg.InitCommonFlags()
	cobraAdaptor := cmds.NewCobraAdaptor(comflag)
	cobraAdaptor.Cmd().Execute()
}

// todo  完成文件输入结构.
