package client

import (
	"context"
)

// PluginRemove removes a plugin
func (cli *Client) PluginRemove(ctx context.Context, name string) error {
	resp, err := cli.delete(ctx, "/plugins/"+name, nil, nil)
	ensureReaderClosed(resp)
	return err
}
