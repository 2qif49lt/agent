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

func checkSystem(broot bool) error {
	// check OS version for compatibility, ensure running in global zone
	var err error

	if broot == true {
		var id C.zoneid_t

		if id, err = C.getzoneid(); err != nil {
			return fmt.Errorf("Exiting. Error getting zone id: %+v", err)
		}
		if int(id) != 0 {
			return fmt.Errorf("Exiting because the agent daemon is not running in the global zone")
		}
	}

	return err
}
