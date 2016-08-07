package client

import (
	"fmt"
	"net/url"

	"golang.org/x/net/context"
)

// Ping returns the agentd server Pong response message.
func (cli *Client) Ping(ctx context.Context, ping string) (string, error) {
	v := url.Values{}
	v.Set("msg", ping)
	serverResp, err := cli.get(ctx, "/ping", v, nil)
	if err != nil {
		return "", err
	}
	defer ensureReaderClosed(serverResp)

	return serverResp, nil
}
