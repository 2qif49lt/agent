// +build linux freebsd darwin

package daemon

import (
	"fmt"
	"io/ioutil"
	"os"
	"runtime/debug"
	"strconv"
	"strings"

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

// verifyDaemonSettings performs validation of daemon config struct
func verifyDaemonSettings(config *Config) error {

	return nil
}

// checkSystem validates platform-specific requirements
func checkSystem() error {
	if os.Geteuid() != 0 {
		return fmt.Errorf("The agent daemon needs to be run as root")
	}
	return checkKernel()
}

// configureMaxThreads sets the Go runtime max threads threshold
// which is 90% of the kernel setting from /proc/sys/kernel/threads-max
func configureMaxThreads(config *Config) error {
	mt, err := ioutil.ReadFile("/proc/sys/kernel/threads-max")
	if err != nil {
		return err
	}
	mtint, err := strconv.Atoi(strings.TrimSpace(string(mt)))
	if err != nil {
		return err
	}
	maxThreads := (mtint / 100) * 90
	debug.SetMaxThreads(maxThreads)
	logrus.Debugf("Golang's threads limit set to %d", maxThreads)
	return nil
}

func (daemon *Daemon) stats() (*types.StatsJSON, error) {
	s := &types.StatsJSON{}
	return s, nil
}
