package client

import (
	"encoding/json"

	"github.com/2qif49lt/agent/api/types"
	"golang.org/x/net/context"
)

// PluginInspect inspects an existing plugin
func (cli *Client) PluginInspect(ctx context.Context, name string) (*types.Plugin, error) {
	var p types.Plugin
	resp, err := cli.get(ctx, "/plugins/"+name, nil, nil)
	if err != nil {
		return nil, err
	}
	err = json.NewDecoder(resp.body).Decode(&p)
	ensureReaderClosed(resp)
	return &p, err
}
