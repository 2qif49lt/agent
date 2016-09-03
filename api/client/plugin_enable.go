package client

// PluginEnable enables a plugin
func (cli *Client) PluginEnable(name string) error {
	resp, err := cli.post("/plugins/"+name+"/enable", nil, nil, nil)
	ensureReaderClosed(resp)
	return err
}
