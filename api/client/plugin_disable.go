package client

// PluginDisable disables a plugin
func (cli *Client) PluginDisable(name string) error {
	resp, err := cli.post("/plugins/"+name+"/disable", nil, nil, nil)
	ensureReaderClosed(resp)
	return err
}
