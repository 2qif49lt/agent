package daemon

var (
	defaultPidFile = os.Getenv("programdata") + string(os.PathSeparator) + "agentd.pid"
)
