package client

import (
	"context"
	"io"

	"github.com/2qif49lt/agent/api/types"
)

// CommonAPIClient is the common methods between stable and experimental versions of APIClient.
type CommonAPIClient interface {
	ClientVersion() string
	ServerVersion(ctx context.Context) (types.Version, error)
	Events(ctx context.Context, options types.EventsOptions) (io.ReadCloser, error)
	Info(ctx context.Context) (types.Info, error)
	Ping(ctx context.Context, ping string) (string, error)

	PluginList(ctx context.Context) (types.PluginsListResponse, error)
	PluginRemove(ctx context.Context, name string) error
	PluginEnable(ctx context.Context, name string) error
	PluginDisable(ctx context.Context, name string) error
	PluginInstall(ctx context.Context, name, registryAuth string, acceptAllPermissions, noEnable bool, in io.ReadCloser, out io.Writer) error
	PluginInspect(ctx context.Context, name string) (*types.Plugin, error)
}
