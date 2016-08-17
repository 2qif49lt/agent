package daemon

import (
	"fmt"
	"os"
	"syscall"

	"github.com/2qif49lt/agent/pkg/signal"
	"github.com/2qif49lt/agent/pkg/system"
	"github.com/2qif49lt/logrus"
)

func setupDumpStackTrap() {
	// Windows does not support signals like *nix systems. So instead of
	// trapping on SIGUSR1 to dump stacks, we wait on a Win32 event to be
	// signaled.
	go func() {
		sa := syscall.SecurityAttributes{
			Length: 0,
		}
		ev := "Global\\agent-daemon-" + fmt.Sprint(os.Getpid())
		if h, _ := system.CreateEvent(&sa, false, false, ev); h != 0 {
			logrus.Debugf("Stackdump - waiting signal at %s", ev)
			for {
				syscall.WaitForSingleObject(h, syscall.INFINITE)
				signal.DumpStacks()
			}
		}
	}()
}
