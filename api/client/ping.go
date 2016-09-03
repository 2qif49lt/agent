package client

import (
	"net/url"
)

// Ping returns the agentd server Pong response message.
func (cli *Client) Ping(ping string) (string, error) {
	v := url.Values{}
	v.Set("msg", ping)
	rsp, err := cli.get("/ping", v, nil)
	if err != nil {
		return "", err
	}
	defer ensureReaderClosed(rsp)

	buf := [1024]byte{}
	rsp.body.Read(buf[:])
	return string(buf[:]), nil
}
