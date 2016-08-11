package plugin

import (
	enginetypes "github.com/2qif49lt/agent/api/types"
)

// Backend for Plugin
type Backend interface {
	Disable(name string) error
	Enable(name string) error
	List() ([]enginetypes.Plugin, error)
	Inspect(name string) (enginetypes.Plugin, error)
	Remove(name string) error
	//	Pull(name string, metaHeaders http.Header, authConfig *enginetypes.AuthConfig) (enginetypes.PluginPrivileges, error)
}
