package main

import (
	_ "github.com/2qif49lt/dump"

	"github.com/2qif49lt/agent/cfg"
	"github.com/2qif49lt/agent/cli/cmds"
	"github.com/2qif49lt/agent/utils"
	"github.com/2qif49lt/logrus"
	flag "github.com/2qif49lt/pflag"

	"io/ioutil"
	"math/rand"
	"os"
	"time"
)

var (
	missionfs = flag.NewFlagSet("mission", flag.ContinueOnError)
	fmission  = missionfs.StringP("mission-file", "m", "", "specify the json formated mission file path,the content's field will overwrite the flags")
)

func initProgram() {
	rand.Seed(time.Now().Unix())
	procdir, _ := utils.GetProcAbsDir()
	os.Chdir(procdir)
}

func init() {
	initProgram()
}

// ***todo: complete the memory usage info in osx

func main() {
	args := preCheckMissionArgs()
	comflag := cfg.InitCommonFlags()
	cobraAdaptor := cmds.NewCobraAdaptor(comflag)

	cobraAdaptor.SetArgs(args)
	cobraAdaptor.Execute()
}

func preCheckMissionArgs() []string {
	args := os.Args[1:]
	missionfs.SetOutput(ioutil.Discard)

	err := missionfs.Parse(args)
	if err == nil && *fmission != "" {
		m := NewMission()
		if err = m.Read(*fmission); err != nil {
			logrus.Errorln(err)
			os.Exit(1)
		}
		args, err = m.ToArgs()
		if err != nil {
			logrus.Errorln(err)
			os.Exit(1)
		}
		logrus.WithTryJson(m).Infoln("use mission file specify flags")
		logrus.WithTryJson(args).Infoln("equal command line")
	}
	return args
}
