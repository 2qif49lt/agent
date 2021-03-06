// +build !windows,!solaris

package daemon

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/2qif49lt/agent/daemon/hack"
	"github.com/2qif49lt/agent/pkg/system"
)

// currentUserIsOwner checks whether the current user is the owner of the given
// file.
func currentUserIsOwner(f string) bool {
	if fileInfo, err := system.Stat(f); err == nil && fileInfo != nil {
		if int(fileInfo.UID()) == os.Getuid() {
			return true
		}
	}
	return false
}

// setDefaultUmask sets the umask to 0022 to avoid problems
// caused by custom umask
func setDefaultUmask() error {
	desiredUmask := 0022
	syscall.Umask(desiredUmask)
	if umask := syscall.Umask(desiredUmask); umask != desiredUmask {
		return fmt.Errorf("failed to set umask: expected %#o, got %#o", desiredUmask, umask)
	}

	return nil
}

// setupConfigReloadTrap configures the USR2 signal to reload the configuration.
func (cli *DaemonCli) setupConfigReloadTrap() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP)
	go func() {
		for range c {
			cli.reloadConfig()
		}
	}()
}

func wrapListeners(proto string, ls net.Listener) net.Listener {
	if proto == "unix" || proto == "fd" {
		ls = &hack.MalformedHostHeaderOverride{ls}
	}
	return ls
}

// notifyShutdown is called after the daemon shuts down but before the process exits.
func notifyShutdown(err error) {
}
