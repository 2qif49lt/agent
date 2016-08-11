// +build experimental

package client

import (
	"io"

	"github.com/2qif49lt/agent/api/types"
	"golang.org/x/net/context"
)

// APIClient is an interface that clients that talk with a docker server must implement.
type APIClient interface {
	CommonAPIClient
	Experimental(ctx context.Context) error
}

// Ensure that Client always implements APIClient.
var _ APIClient = &Client{}
