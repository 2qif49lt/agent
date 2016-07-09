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

	// DefaultUniqueIdFile 和证书同时生成的用于表示服务器唯一性,合法性.
	DefaultUniqueAgentIdFile = "agentid"

	// DefaultKeyFile is the default filename for the key pem file
	DefaultKeyFile = "agent_key.pem"
	// DefaultCertFile is the default filename for the cert pem file
	DefaultCertFile = "agent_cert.pem"
	// TLSVerifyKey is the default flag name for the tls verification option

	DefaultPublicKeyFile = "server_pub.pem" // 用于检查调用者参数签名
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

// ConfigDir returns the directory the configuration file is stored in
func ConfigDir() string {
	return configDir
}

func CertPath() string {
	return certPath
}

// SetConfigDir sets the directory the configuration file is stored in
func SetConfigDir(dir string) {
	configDir = dir
}
func SetCertPath(dir string) {
	certPath = dir
}

// Newcfgfile initializes an empty configuration file for the given file path 'fp'
func Newcfgfile(fp string) *cfgfile.ConfigFile {
	return &cfgfile.ConfigFile{
		SrvName:  "Agentd",
		Port:     1688,
		Loglvl:   1,
		Filename: fp,
	}
}

// Load reads the configuration files
func Load() (*cfgfile.ConfigFile, error) {

	cfgfile := Newcfgfile(filepath.Join(configDir, ConfigFileName))

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

// IsLegal return whether agent install properly
func IsLegal(cfg *cfgfile.ConfigFile) error {
	isexist := func(fn string) bool {
		_, err := os.Stat(filepath.Join(certPath(), fn))
		return err == nil || os.IsExist(err)
	}
	if isexist(DefaultCertFile) && isexist(DefaultKeyFile) &&
		isexist(DefaultPublicKeyFile) {

	}
	return fmt.Errorf("cert files not exist")
}
