package cfgfile

import (
	"github.com/2qif49lt/toml"

	"fmt"
	"io"
	"os"
	"path/filepath"
)

// ConfigFile ~/.agent/config.toml file info,for agentd.

type ConfigFile struct {
	SrvName string `toml:"srvname,omitempty"` // default Agentd
	Srv     string `toml:"srv,omitempty"`     // if ip is empty, listen on 127.0.0.1 only
	Loglvl  string `toml:"loglvl,omitempty"`  // default InfoLevel
	Master  struct {
		Srvs []string `toml:"services,omitempty"`
	}
	Filename string `toml:"-"` // Note: for internal use only

}

// LoadFromReader reads the configuration data given
func (configFile *ConfigFile) LoadFromReader(configData io.Reader) error {
	if err := toml.NewDecoder(configData).Decode(configFile); err != nil {
		return err
	}

	return nil
}

func (configFile *ConfigFile) SaveToWriter(writer io.Writer) error {

	data, err := toml.Marshal(*configFile)
	if err != nil {
		return err
	}
	_, err = writer.Write(data)
	return err
}

// Save encodes and writes out all information
func (configFile *ConfigFile) Save() error {
	if configFile.Filename == "" {
		return fmt.Errorf("Can't save config with empty filename")
	}

	if err := os.MkdirAll(filepath.Dir(configFile.Filename), 0700); err != nil {
		return err
	}
	f, err := os.OpenFile(configFile.Filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer f.Close()
	return configFile.SaveToWriter(f)
}
