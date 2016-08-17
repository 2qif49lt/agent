package daemon

import (
	// "github.com/Microsoft/hcsshim"
	"github.com/2qif49lt/agent/pkg/system"
)

const (
	platformSupported = true
)

// checkSystem validates platform-specific requirements
func checkSystem(broot bool) error {

	return nil
	// Validate the OS version. Note that docker.exe must be manifested for this
	// call to return the correct version.
	osv := system.GetOSVersion()
	if osv.MajorVersion < 10 {
		return fmt.Errorf("This version of Windows does not support the agent daemon")
	}
	if osv.Build < 14300 {
		return fmt.Errorf("The Windows daemon requires Windows Server 2016 Technical Preview 5 build 14300 or later")
	}
	return nil
}
