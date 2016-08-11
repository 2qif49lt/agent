package client

import (
	"net/http"
	"runtime"

	"github.com/2qif49lt/agent/api"
	apiclient "github.com/2qif49lt/agent/api/client"
	"github.com/2qif49lt/agent/cfg"
	"github.com/2qif49lt/agent/pkg/connections/sockets"
	"github.com/2qif49lt/agent/pkg/connections/tlsconfig"
	"github.com/2qif49lt/agent/pkg/opts"
	"github.com/2qif49lt/agent/version"
)

// Cli represents the agent command line client.
// Instances of the client can be returned from NewDockerCli.
type AgentCli struct {
	// initializing closure
	init   func() error
	common *cfg.CommonFlags

	// api client
	client apiclient.Client
}

// Initialize calls the init function that will setup the configuration for the client
// such as the TLS, tcp and other parameters used to run the client.
func (cli *AgentCli) Initialize() error {
	if cli.init == nil {
		return nil
	}
	return cli.init()
}

// Client returns the APIClient
func (cli *AgentCli) Client() apiclient.Client {
	return cli.client
}

// NewAgentCli returns a agent client instance with config.
// The key file, protocol (i.e. unix) and address are passed in as strings, along with the tls.Config.
// If the tls.Config is set the client scheme will be set to https.
func NewAgentCli(com *cfg.CommonFlags) *AgentCli {
	cli := &Cli{
		common: com,
	}

	cli.init = func() error {
		client, err := NewAPIClientFromFlags(com)
		if err != nil {
			return err
		}

		cli.client = client

		return nil
	}

	return cli
}

// NewAPIClientFromFlags creates a new APIClient from command line flags
func NewAPIClientFromFlags(com *cfg.CommonFlags) (apiclient.Client, error) {
	host, err = opts.ParseHost(!com.NoTLS, com.Host)

	customHeaders := map[string]string{}
	customHeaders["User-Agent"] = clientUserAgent()

	verStr := api.DefaultVersion
	if tmpStr := os.Getenv("AGENT_API_VERSION"); tmpStr != "" {
		verStr = tmpStr
	}

	httpClient, err := newHTTPClient(host, com.TLSOptions)
	if err != nil {
		return &client.Client{}, err
	}

	return client.NewClient(host, verStr, httpClient, customHeaders)
}

func newHTTPClient(host string, tlsOptions *tlsconfig.Options) (*http.Client, error) {
	if tlsOptions == nil {
		// let the api client configure the default transport.
		return nil, nil
	}

	config, err := tlsconfig.Client(*tlsOptions)
	if err != nil {
		return nil, err
	}
	tr := &http.Transport{
		TLSClientConfig: config,
	}
	proto, addr, _, err := client.ParseHost(host)
	if err != nil {
		return nil, err
	}

	sockets.ConfigureTransport(tr, proto, addr)

	return &http.Client{
		Transport: tr,
	}, nil
}

func clientUserAgent() string {
	return "Agent-Client/" + version.CLI_VERSION + " (" + runtime.GOOS + ")"
}
