// +build linux freebsd

package daemon

import (
	"fmt"
	"os"
	"runtime/debug"
	"strconv"
	"strings"
	"time"

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

// adaptContainerSettings is called during container creation to modify any
// settings necessary in the HostConfig structure.
func (daemon *Daemon) adaptContainerSettings(hostConfig *containertypes.HostConfig, adjustCPUShares bool) error {
	if adjustCPUShares && hostConfig.CPUShares > 0 {
		// Handle unsupported CPUShares
		if hostConfig.CPUShares < linuxMinCPUShares {
			logrus.Warnf("Changing requested CPUShares of %d to minimum allowed of %d", hostConfig.CPUShares, linuxMinCPUShares)
			hostConfig.CPUShares = linuxMinCPUShares
		} else if hostConfig.CPUShares > linuxMaxCPUShares {
			logrus.Warnf("Changing requested CPUShares of %d to maximum allowed of %d", hostConfig.CPUShares, linuxMaxCPUShares)
			hostConfig.CPUShares = linuxMaxCPUShares
		}
	}
	if hostConfig.Memory > 0 && hostConfig.MemorySwap == 0 {
		// By default, MemorySwap is set to twice the size of Memory.
		hostConfig.MemorySwap = hostConfig.Memory * 2
	}
	if hostConfig.ShmSize == 0 {
		hostConfig.ShmSize = container.DefaultSHMSize
	}
	var err error
	if hostConfig.SecurityOpt == nil {
		hostConfig.SecurityOpt, err = daemon.generateSecurityOpt(hostConfig.IpcMode, hostConfig.PidMode, hostConfig.Privileged)
		if err != nil {
			return err
		}
	}
	if hostConfig.MemorySwappiness == nil {
		defaultSwappiness := int64(-1)
		hostConfig.MemorySwappiness = &defaultSwappiness
	}
	if hostConfig.OomKillDisable == nil {
		defaultOomKillDisable := false
		hostConfig.OomKillDisable = &defaultOomKillDisable
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
	if !c.IsRunning() {
		return nil, fmt.Errorf(`agent daemon is not running`)
	}
	s := &types.StatsJSON{}
	s.Read = time.Unix(int64(stats.Timestamp), 0)
	return s, nil
}
