package main

import (
	_ "github.com/2qif49lt/dump"

	"github.com/2qif49lt/agent/cfg"
	"github.com/2qif49lt/agent/cli/cmds"
	"github.com/2qif49lt/agent/utils"
	"github.com/2qif49lt/logrus"
	flag "github.com/2qif49lt/pflag"

	"math/rand"
	"os"
	"time"
)

var (
	fmission = flag.StringP("mission-file", "m", "", "specify the json formated mission file path,the content's field will overwrite the flags")
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
	flag.Parse()
	if *fmission != "" {
		logrus.Infoln(*fmission)
	}

	comflag := cfg.InitCommonFlags()
	cobraAdaptor := cmds.NewCobraAdaptor(comflag)
	cobraAdaptor.Cmd().Execute()
}

// todo  完成文件输入结构.
