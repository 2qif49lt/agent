package client

// PluginRemove removes a plugin
func (cli *Client) PluginRemove(name string) error {
	resp, err := cli.delete("/plugins/"+name, nil, nil)
	ensureReaderClosed(resp)
	return err
}
