// Package cfg implements config's manage

package cfg

import (
	"os"
)

const (
	ConfigFileName = "config.toml"
	ConfigFileDir  = ".agent"
)

var (
	configDirPath = os.Getenv("AGENT_CONFIG_PATH")
)

func init() {
	if configDir == "" {
		configDir = filepath.Join(homedir.Get(), configFileDir)
	}

}
