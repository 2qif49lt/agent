// Package cfg implements config's manage

package cfg

import (
	"github.com/2qif49lt/agent/cfg/cfgfile"
	"github.com/2qif49lt/agent/utils"
	"github.com/2qif49lt/logrus"
	"path/filepath"

	"io/ioutil"
	"os"
)

const (
	ConfigFileName = "config.toml"

	configFileDir = "config"
	certFileDir   = "cert"

	DefaultAgentdListenPort = 3567 // 默认监听端口

	DefaultUniqueAgentIdFile = "agentid"

	DeafultTlsCaFile  = "ca-cert.pem"
	DefaultTlsKeyFile = "tls-key.pem" // 作为服务器和客户端时所需的证书不一样.
	DefultTlsCertFile = "tls-cert.pem"

	// DefaultSignPubFile 参数签名
	DefaultRsaSignFile = "rsa-pub.pem" // 用于检查调用者参数签名,agentd使用公钥,agent使用密钥.
)

var Conf *cfgfile.ConfigFile

var (
	configDir = os.Getenv("AGENT_CONFIG_PATH") // agent配置文件目录
	certPath  = os.Getenv("AGENT_CERT_PATH")   // 证书所在目录,默认在.agent
)

func init() {
	prgpath, _ := utils.GetProcAbsDir()
	if configDir == "" {
		configDir = filepath.Join(prgpath, configFileDir)
	}
	if certPath == "" {
		certPath = filepath.Join(prgpath, certFileDir)
	}
}

// GetConfigDir returns the directory the configuration file is stored in
func GetConfigDir() string {
	return configDir
}

func GetCertPath() string {
	return certPath
}

var changeConfigDir = false
var changeCertPahth = false

// SetConfigDir sets the directory the configuration file is stored in
func SetConfigDir(dir string) {
	changeConfigDir = true
	configDir = dir
}
func SetCertPath(dir string) {
	changeCertPahth = true
	certPath = dir
}

// newcfgfile initializes an empty configuration file for the given file path 'fp'
func newcfgfile(fp string) *cfgfile.ConfigFile {
	return &cfgfile.ConfigFile{
		SrvName:  "agentd",
		Filename: fp,
	}
}

func InitConf() error {
	_, err := load()
	return err
}

// Load reads the configuration files
func load() (*cfgfile.ConfigFile, error) {
	if changeConfigDir {
		logrus.Infoln("use config:", filepath.Join(configDir, ConfigFileName))
	}
	conf := newcfgfile(filepath.Join(configDir, ConfigFileName))
	err := conf.Load()
	if err != nil && err != cfgfile.ErrConfigFileMiss {
		return nil, err
	}
	prgpath, _ := utils.GetProcAbsDir()
	tmp, err := ioutil.ReadFile(filepath.Join(prgpath, DefaultUniqueAgentIdFile))
	conf.Agentid = string(tmp)

	if err == nil {
		Conf = conf
	}

	return conf, err
}
