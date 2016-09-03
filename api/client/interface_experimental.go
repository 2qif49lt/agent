// +build experimental

package client

// APIClient is an interface that clients that talk with a docker server must implement.
type APIClient interface {
	CommonAPIClient
	Experimental(ctx context.Context) error
}

// Ensure that Client always implements APIClient.
var _ APIClient = &Client{}

// Info returns information about the agentd server.
func (cli *Client) Experimental() error {
	return nil
}
