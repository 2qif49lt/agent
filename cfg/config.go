// Package cfg implements config's manage

package cfg

import (
	"github.com/2qif49lt/agent/cfg/cfgfile"
	"github.com/2qif49lt/agent/pkg/homedir"
	"path/filepath"

	"fmt"
	"os"
)

const (
	ConfigFileName = "config.toml"
	configFileDir  = ".agent"

	DefaultAgentdListenPort = 1698 // 默认监听端口

	// DefaultUniqueIdFile 和证书同时生成的用于表示服务器唯一性,合法性.
	DefaultUniqueAgentIdFile = "agentid"

	// DeafultTlsCaFile 用于判断连接的agentd的证书是否合法
	DeafultTlsCaFile = "ca.pem"
	// DefaultTlsKeyFile agentd的tls链接key
	DefaultTlsKeyFile = "agent_key.pem"
	// DeafultTlsCertFile agentd的证书.
	DefultTlsCertFile = "agent_cert.pem"

	// DefaultSignPubFile 参数签名
	DefaultRsaSignPubFile = "server_pub.pem" // 用于agentd检查调用者参数签名
	DefaultRsaSignPriFile = "server_key.pem" // 用于客户端签名
)

var (
	configDir = os.Getenv("AGENT_CONFIG_PATH") // agent配置文件目录
	certPath  = os.Getenv("AGENT_CERT_PATH")   // 证书所在目录,默认在.agent
)

func init() {
	if configDir == "" {
		configDir = filepath.Join(homedir.Get(), configFileDir)
	}
	if certPath == "" {
		certPath = configDir
	}
}

// GetConfigDir returns the directory the configuration file is stored in
func GetConfigDir() string {
	return configDir
}

func GetCertPath() string {
	return certPath
}

// SetConfigDir sets the directory the configuration file is stored in
func SetConfigDir(dir string) {
	configDir = dir
}
func SetCertPath(dir string) {
	certPath = dir
}

// newcfgfile initializes an empty configuration file for the given file path 'fp'
func newcfgfile(fp string) *cfgfile.ConfigFile {
	return &cfgfile.ConfigFile{
		SrvName:  "Agentd",
		Filename: fp,
	}
}

// Load reads the configuration files
func Load() (*cfgfile.ConfigFile, error) {

	cfgfile := newcfgfile(filepath.Join(GetConfigDir(), ConfigFileName))

	if _, err := os.Stat(cfgfile.Filename); err == nil {
		file, err := os.Open(cfgfile.Filename)
		if err != nil {
			return nil, fmt.Errorf("%s - %v", cfgfile.Filename, err)
		}
		defer file.Close()

		err = cfgfile.LoadFromReader(file)
		if err != nil {
			err = fmt.Errorf("%s - %v", cfgfile.Filename, err)
			return nil, err
		}
		return cfgfile, nil
	} else {
		return nil, fmt.Errorf("%s - %v", cfgfile.Filename, err)
	}

}

// IsTlsLegal return whether agentd install properly
func IsSrvTlsLegal(cfg *cfgfile.ConfigFile) error {
	isexist := func(fn string) bool {
		_, err := os.Stat(filepath.Join(GetCertPath(), fn))
		return err == nil || os.IsExist(err)
	}
	if isexist(DefaultTlsKeyFile) && isexist(DefultTlsCertFile) {

	}
	return fmt.Errorf("cert files not exist")
}

func IsCliTlsLegal() error {
	return nil
}
