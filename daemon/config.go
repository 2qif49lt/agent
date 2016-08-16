package daemon

import (
	"sync"

	"github.com/2qif49lt/logrus"
	flag "github.com/2qif49lt/pflag"
)

// Config defines the configuration of a agent daemon.
type Config struct {
	Root              bool // 是否以root启动
	OSMaxThreadNum    int  // 程序最多占有多少线程数. 单位CPU个数.比如2 表示2倍CPU数量的线程
	AgentID           string
	CorsHeaders       string
	NoTLSClientVerify bool // 是否要求客户端验证
	SocketGroup       string
	reloadLock        sync.Mutex
}

// ReloadConfiguration reads the configuration in the host and reloads the daemon and server.
func ReloadConfiguration(configFile string, flags *flag.FlagSet, reload func(*Config)) error {
	logrus.Infof("Got signal to reload configuration, reloading from: %s", configFile)

	return nil
}

// InstallFlags adds command-line options to the top-level flag parser for
// the current process.
// Subsequent calls to `flag.Parse` will populate config with values parsed
// from the command-line.
func (config *Config) InstallFlags(cmd *flag.FlagSet) {
	cmd.BoolVarP(&config.Root, "root", "r", true, "run agent as root")
	cmd.IntVarP(&config.OSMaxThreadNum, "thread-num", "t", 0, "Set the maximum OS threads used by agentd,unit is the num of logical cpu.")

	cmd.StringVar(&config.AgentID, "agent-id", "", "Set CORS headers in the remote API")
	cmd.StringVar(&config.CorsHeaders, "api-cors-header", "", "Set CORS headers in the remote API")
	cmd.BoolVar(&config.NoTLSClientVerify, "noverify", false, "DO NOT verify client certificate")
	cmd.StringVarP(&config.SocketGroup, "group", "g", "agentd", "Group name for the unix socket")
}
