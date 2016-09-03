package client

import (
	"encoding/json"

	"github.com/2qif49lt/agent/api/types"
)

// PluginInspect inspects an existing plugin
func (cli *Client) PluginInspect(name string) (*types.Plugin, error) {
	var p types.Plugin
	resp, err := cli.get("/plugins/"+name, nil, nil)
	if err != nil {
		return nil, err
	}
	err = json.NewDecoder(resp.body).Decode(&p)
	ensureReaderClosed(resp)
	return &p, err
}
