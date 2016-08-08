package system

import (
	"fmt"
	"os"
	"runtime"
	"time"

	"golang.org/x/net/context"

	"github.com/2qif49lt/agent/api/types"
	"github.com/2qif49lt/agent/cli"
	"github.com/2qif49lt/agent/client"
	"github.com/2qif49lt/agent/utils"
	"github.com/2qif49lt/agent/utils/templates"
	"github.com/2qif49lt/agent/version"
	"github.com/2qif49lt/cobra"
)

var versionTemplate = `Client:
 Version:      {{.Client.Version}}
 API version:  {{.Client.APIVersion}}
 Go version:   {{.Client.GoVersion}}
 Built:        {{.Client.BuildTime}}
 OS/Arch:      {{.Client.Os}}/{{.Client.Arch}}{{if .Client.Experimental}}
 Experimental: {{.Client.Experimental}}{{end}}{{if .ServerOK}}

Server:
 Version:      {{.Server.Version}}
 API version:  {{.Server.APIVersion}}
 Go version:   {{.Server.GoVersion}}
 Built:        {{.Server.BuildTime}}
 OS/Arch:      {{.Server.Os}}/{{.Server.Arch}}{{if .Server.Experimental}}
 Experimental: {{.Server.Experimental}}{{end}}{{end}}`

type versionOptions struct {
	format string
}

// NewVersionCommand creats a new cobra.Command for `agent version`
func NewVersionCommand(agentCli *client.AgentCli) *cobra.Command {
	var opts versionOptions

	cmd := &cobra.Command{
		Use:   "version [OPTIONS]",
		Short: "Show the agent version information",
		Args:  cli.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runVersion(agentCli, &opts)
		},
	}

	flags := cmd.Flags()

	flags.StringVarP(&opts.format, "format", "f", "", "Format the output using the given go template")

	return cmd
}

func runVersion(agentCli *client.AgentCli, opts *versionOptions) error {
	ctx := context.Background()

	templateFormat := versionTemplate
	if opts.format != "" {
		templateFormat = opts.format
	}

	tmpl, err := templates.Parse(templateFormat)
	if err != nil {
		return cli.StatusError{StatusCode: 64,
			Status: "Template parsing error: " + err.Error()}
	}

	vd := types.VersionResponse{
		Client: &types.Version{
			Version:      version.CLI_VERSION,
			APIVersion:   agentCli.Client().ClientVersion(),
			GoVersion:    runtime.Version(),
			BuildTime:    version.BUILDTIME,
			Os:           runtime.GOOS,
			Arch:         runtime.GOARCH,
			Experimental: utils.ExperimentalBuild(),
		},
	}

	serverVersion, err := agentCli.Client().ServerVersion(ctx)
	if err == nil {
		vd.Server = &serverVersion
	}

	// first we need to make BuildTime more human friendly
	t, errTime := time.Parse(time.RFC3339Nano, vd.Client.BuildTime)
	if errTime == nil {
		vd.Client.BuildTime = t.Format(time.ANSIC)
	}

	if vd.ServerOK() {
		t, errTime = time.Parse(time.RFC3339Nano, vd.Server.BuildTime)
		if errTime == nil {
			vd.Server.BuildTime = t.Format(time.ANSIC)
		}
	}

	if err2 := tmpl.Execute(os.Stdout, vd); err2 != nil && err == nil {
		err = err2
	}
	fmt.Printf("\n")
	return err
}
