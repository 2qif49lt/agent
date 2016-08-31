package daemon

import (
	"os"
	"runtime"
	"time"

	"github.com/2qif49lt/agent/api"
	"github.com/2qif49lt/agent/api/types"
	"github.com/2qif49lt/agent/pkg/connections/sockets"
	"github.com/2qif49lt/agent/pkg/fileutils"
	"github.com/2qif49lt/agent/pkg/parsers/kernel"
	"github.com/2qif49lt/agent/pkg/parsers/operatingsystem"
	"github.com/2qif49lt/agent/pkg/platform"
	"github.com/2qif49lt/agent/pkg/system"
	"github.com/2qif49lt/agent/utils"
	"github.com/2qif49lt/logrus"
)

// SystemInfo returns information about the host server the daemon is running on.
func (daemon *Daemon) SystemInfo() (*types.Info, error) {
	kernelVersion := "<unknown>"
	if kv, err := kernel.GetKernelVersion(); err != nil {
		logrus.Warnf("Could not get kernel version: %v", err)
	} else {
		kernelVersion = kv.String()
	}

	operatingSystem := "<unknown>"
	if s, err := operatingsystem.GetOperatingSystem(); err != nil {
		logrus.Warnf("Could not get operating system name: %v", err)
	} else {
		operatingSystem = s
	}

	meminfo, err := system.ReadMemInfo()
	if err != nil {
		logrus.Errorf("Could not read system memory info: %v", err)
		return nil, err
	}

	v := &types.Info{
		ID:                daemon.configStore.AgentID,
		Debug:             utils.IsDebugEnabled(),
		NFd:               fileutils.GetTotalUsedFds(),
		NGoroutines:       runtime.NumGoroutine(),
		SystemTime:        time.Now().Format("2006-01-02 15:04:05"),
		NEventsListener:   daemon.EventsService.SubscribersCount(),
		KernelVersion:     kernelVersion,
		OperatingSystem:   operatingSystem,
		OSType:            platform.OSType,
		Architecture:      platform.Architecture,
		NCPU:              runtime.NumCPU(),
		MemTotal:          meminfo.MemTotal,
		ExperimentalBuild: utils.ExperimentalBuild(),
		ServerVersion:     api.API_VERSION,
		HTTPProxy:         sockets.GetProxyEnv("http_proxy"),
		HTTPSProxy:        sockets.GetProxyEnv("https_proxy"),
		NoProxy:           sockets.GetProxyEnv("no_proxy"),
	}

	hostname := ""
	if hn, err := os.Hostname(); err != nil {
		logrus.Warnf("Could not get hostname: %v", err)
	} else {
		hostname = hn
	}
	v.Name = hostname

	return v, nil
}

// SystemVersion returns version information about the daemon.
func (daemon *Daemon) SystemVersion() types.Version {
	v := types.Version{
		Version:      api.API_VERSION,
		GoVersion:    runtime.Version(),
		Os:           runtime.GOOS,
		Arch:         runtime.GOARCH,
		BuildTime:    api.BUILDTIME,
		Experimental: utils.ExperimentalBuild(),
	}

	kernelVersion := "<unknown>"
	if kv, err := kernel.GetKernelVersion(); err != nil {
		logrus.Warnf("Could not get kernel version: %v", err)
	} else {
		kernelVersion = kv.String()
	}
	v.KernelVersion = kernelVersion

	return v
}
