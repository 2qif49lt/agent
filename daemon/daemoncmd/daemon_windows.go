package main

import (
	"fmt"
	"os"
	"syscall"

	"github.com/2qif49lt/agent/pkg/system"
	"github.com/2qif49lt/logrus"
)

// currentUserIsOwner checks whether the current user is the owner of the given
// file.
func currentUserIsOwner(f string) bool {
	return false
}

// setDefaultUmask doesn't do anything on windows
func setDefaultUmask() error {
	return nil
}

// setupConfigReloadTrap configures a Win32 event to reload the configuration.
func (cli *DaemonCli) setupConfigReloadTrap() {
	go func() {
		sa := syscall.SecurityAttributes{
			Length: 0,
		}
		ev := "Global\\agent-daemon-config-" + fmt.Sprint(os.Getpid())
		if h, _ := system.CreateEvent(&sa, false, false, ev); h != 0 {
			logrus.Debugf("Config reload - waiting signal at %s", ev)
			for {
				syscall.WaitForSingleObject(h, syscall.INFINITE)
				cli.reloadConfig()
			}
		}
	}()
}

func wrapListeners(proto string, ls net.Listener) net.Listener {
	return ls
}
