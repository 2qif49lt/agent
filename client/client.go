package client

import (
	"net/http"
	"os"
	"runtime"

	"github.com/2qif49lt/agent/api"
	apiclient "github.com/2qif49lt/agent/api/client"
	"github.com/2qif49lt/agent/cfg"
	"github.com/2qif49lt/agent/pkg/connections/sockets"
	"github.com/2qif49lt/agent/pkg/connections/tlsconfig"
	"github.com/2qif49lt/agent/pkg/opts"
	flag "github.com/2qif49lt/pflag"
)

type Config struct {
	*cfg.CommonFlags
	// 操作命令的JSON格式
	Mission string
}

// AgentCli represents the agent command line client.
// Instances of the client can be returned from NewDockerCli.
type AgentCli struct {
	*Config

	// initializing closure
	init func() error
	// api client
	client apiclient.APIClient
}

func (config *Config) installFlags(flags *flag.FlagSet) {
	//	flags.StringVarP(&config.Mission, "mission-file", "m", "", "specify the json formated mission file path,the content's field will overwrite the flags")
}
func (cli *AgentCli) InitFlags(cmd *flag.FlagSet) {
	cli.Config.installFlags(cmd)
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
func (cli *AgentCli) Client() apiclient.APIClient {
	return cli.client
}

// NewAgentCli returns a agent client instance with config.
// The key file, protocol (i.e. unix) and address are passed in as strings, along with the tls.Config.
// If the tls.Config is set the client scheme will be set to https.
func NewAgentCli(com *cfg.CommonFlags) *AgentCli {
	cli := &AgentCli{
		Config: &Config{
			CommonFlags: com,
		},
	}

	cli.init = func() error {
		client, err := NewAPIClient(cli.Config)
		if err != nil {
			return err
		}

		cli.client = client

		return nil
	}

	return cli
}

// NewAPIClient creates a new APIClient
func NewAPIClient(clicfg *Config) (apiclient.APIClient, error) {
	com := clicfg.CommonFlags
	host, err := opts.ParseHost(com.Host)
	if err != nil {
		return nil, err
	}

	customHeaders := map[string]string{}
	customHeaders["User-Agent"] = clientUserAgent()

	verStr := api.API_VERSION
	if tmpStr := os.Getenv("AGENT_API_VERSION"); tmpStr != "" {
		verStr = tmpStr
	}
	customHeaders["version"] = verStr

	httpClient, err := newHTTPClient(host, com.TLSOptions)
	if err != nil {
		return nil, err
	}

	return apiclient.NewClient(host, verStr, httpClient, customHeaders)
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
	proto, addr, _, err := apiclient.ParseHost(host)
	if err != nil {
		return nil, err
	}

	sockets.ConfigureTransport(tr, proto, addr)

	return &http.Client{
		Transport: tr,
	}, nil
}

func clientUserAgent() string {
	return "Agent-Client/" + runtime.GOOS
}
