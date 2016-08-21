package daemon

import (
	"sync"

	"github.com/2qif49lt/logrus"
	flag "github.com/2qif49lt/pflag"
)

// Config defines the configuration of a agent daemon.
type Config struct {
	SrvName           string // 服务名
	Root              bool   // 是否以root启动
	OSMaxThreadNum    int    // 程序最多占有多少线程数. 单位CPU个数.比如2 表示2倍CPU数量的线程
	AgentID           string
	CorsHeaders       string
	NoTLSClientVerify bool // 是否要求客户端验证
	SocketGroup       string
	CertExtenAuth     bool // 对证书1.2.3.4字段进行检查

	reloadLock sync.Mutex
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
func (config *Config) InstallFlags(flags *flag.FlagSet) {
	flags.StringVarP(&config.SrvName, "name", "n", "", "指定服务名,若空则使用配置文件内值,若无配置则默认")
	flags.BoolVarP(&config.Root, "root", "r", true, "Run agent as root")
	flags.IntVarP(&config.OSMaxThreadNum, "thread-num", "t", 0, "Set the maximum OS threads used by agentd,unit is the num of logical cpu.")
	flags.StringVar(&config.AgentID, "agent-id", "", "Specify agent id")
	flags.StringVar(&config.CorsHeaders, "api-cors-header", "", "Set CORS headers in the remote API")
	flags.BoolVar(&config.NoTLSClientVerify, "noverify", false, "DO NOT verify client certificate")
	flags.StringVarP(&config.SocketGroup, "group", "g", "agentd", "Group name for the unix socket")
	flags.BoolVar(&config.CertExtenAuth, "cert-exten-auth", false, "check the extenion field for routing")

}
