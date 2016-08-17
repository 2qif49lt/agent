// +build linux freebsd darwin

package daemon

import (
	"fmt"
	"os"

	"github.com/2qif49lt/agent/api/types"
	"github.com/2qif49lt/agent/pkg/parsers/kernel"
	"github.com/2qif49lt/logrus"
)

const (
	platformSupported = true
)

func checkKernelVersion(k, major, minor int) bool {
	if v, err := kernel.GetKernelVersion(); err != nil {
		logrus.Warnf("error getting kernel version: %s", err)
	} else {
		if kernel.CompareKernelVersion(*v, kernel.VersionInfo{Kernel: k, Major: major, Minor: minor}) < 0 {
			return false
		}
	}
	return true
}

func checkKernel() error {

	if !checkKernelVersion(3, 10, 0) {
		v, _ := kernel.GetKernelVersion()
		if os.Getenv("AGENT_NOWARN_KERNEL_VERSION") == "" {
			logrus.Fatalf("Your Linux kernel version %s is not supported for running agent. Please upgrade your kernel to 3.10.0 or newer.", v.String())
		}
	}
	return nil
}

// checkSystem validates platform-specific requirements
func checkSystem(broot bool) error {
	if broot == true {
		if os.Geteuid() != 0 {
			return fmt.Errorf("The agent daemon needs to be run as root")
		}
	}

	return nil
}

func (daemon *Daemon) stats() (*types.StatsJSON, error) {
	s := &types.StatsJSON{}
	return s, nil
}
