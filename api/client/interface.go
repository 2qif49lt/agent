package client

import (
	"io"

	"github.com/2qif49lt/agent/api/types"
)

// CommonAPIClient is the common methods between stable and experimental versions of APIClient.
type CommonAPIClient interface {
	ClientVersion() string
	ServerVersion() (types.Version, error)
	Events(options types.EventsOptions) (io.ReadCloser, error)
	Info() (types.Info, error)
	Ping(ping string) (string, error)

	PluginList() (types.PluginsListResponse, error)
	PluginRemove(name string) error
	PluginEnable(name string) error
	PluginDisable(name string) error
	PluginInstall(name, registryAuth string, acceptAllPermissions, noEnable bool, in io.ReadCloser, out io.Writer) error
	PluginInspect(name string) (*types.Plugin, error)
}
