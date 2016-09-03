package client

import (
	"encoding/json"

	"github.com/2qif49lt/agent/api/types"
)

// PluginList returns the installed plugins
func (cli *Client) PluginList() (types.PluginsListResponse, error) {
	var plugins types.PluginsListResponse
	resp, err := cli.get("/plugins", nil, nil)
	if err != nil {
		return plugins, err
	}

	err = json.NewDecoder(resp.body).Decode(&plugins)
	ensureReaderClosed(resp)
	return plugins, err
}
