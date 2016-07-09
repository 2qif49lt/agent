package cfg

import (
	"path/filepath"
	"testing"
)

func TestNormalUser(t *testing.T) {
	cfg, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Filename != filepath.Join(configDir, ConfigFileName) {
		t.Fatal(cfg.Filename, configDir, ConfigFileName)
	}
	if cfg.Loglvl == 1 {
		cfg.Loglvl = 2
	} else {
		cfg.Loglvl = 1
	}
	cfg.Etcd.Srvs = nil
	t.Log(cfg)
	err = cfg.Save()
	if err != nil {
		t.Fatal(err)
	}
}
