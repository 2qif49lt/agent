// +build solaris,cgo

package daemon

import (
	"fmt"

	"github.com/2qif49lt/agent/api/types"
	"github.com/2qif49lt/agent/pkg/parsers/kernel"
)

//#include <zone.h>
import "C"

const (
	platformSupported = true
)

func checkKernel() error {
	// solaris can rely upon checkSystem() below, we don't skew kernel versions
	return nil
}

// verifyDaemonSettings performs validation of daemon config struct
func verifyDaemonSettings(config *Config) error {
	// checkSystem validates platform-specific requirements
	return nil
}

func checkSystem() error {
	// check OS version for compatibility, ensure running in global zone
	var err error
	var id C.zoneid_t

	if id, err = C.getzoneid(); err != nil {
		return fmt.Errorf("Exiting. Error getting zone id: %+v", err)
	}
	if int(id) != 0 {
		return fmt.Errorf("Exiting because the Docker daemon is not running in the global zone")
	}

	v, err := kernel.GetKernelVersion()
	if kernel.CompareKernelVersion(*v, kernel.VersionInfo{Kernel: 5, Major: 12, Minor: 0}) < 0 {
		return fmt.Errorf("Your Solaris kernel version: %s doesn't support Docker. Please upgrade to 5.12.0", v.String())
	}
	return err
}
