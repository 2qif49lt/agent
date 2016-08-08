package daemon

import (
	"os"
)

var (
	defaultPidFile = os.Getenv("programdata") + string(os.PathSeparator) + "agentd.pid"
)
