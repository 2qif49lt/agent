package system

import (
	"context"
	"fmt"
	"os"
	"runtime"

	"github.com/2qif49lt/agent/api"
	"github.com/2qif49lt/agent/api/types"
	"github.com/2qif49lt/agent/cli"
	"github.com/2qif49lt/agent/client"
	"github.com/2qif49lt/agent/pkg/parsers/kernel"
	"github.com/2qif49lt/agent/utils"
	"github.com/2qif49lt/agent/utils/templates"
	"github.com/2qif49lt/cobra"
	"github.com/2qif49lt/logrus"
)

var versionTemplate = `Client:
 API version:  {{.Client.APIVersion}}
 Go version:   {{.Client.GoVersion}}
 Built:        {{.Client.BuildTime}}
 OS/Arch:      {{.Client.Os}}/{{.Client.Arch}}
 Kernel:       {{.Client.KernelVersion}}{{if .Client.Experimental}}
 Experimental: {{.Client.Experimental}}{{end}}{{if .ServerOK}}

Server:
 API version:  {{.Server.APIVersion}}
 Go version:   {{.Server.GoVersion}}
 Built:        {{.Server.BuildTime}}
 OS/Arch:      {{.Server.Os}}/{{.Server.Arch}}
 Kernel:       {{.Server.KernelVersion}}{{if .Server.Experimental}}
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
		Args:  cli.RequiresMaxArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return agentCli.Initialize()
		},
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
		return err
	}
	kernelVersion := "<unknown>"
	if kv, err := kernel.GetKernelVersion(); err != nil {
		logrus.Warnf("Could not get kernel version: %v", err)
	} else {
		kernelVersion = kv.String()
	}

	vd := types.VersionResponse{
		Client: &types.Version{
			APIVersion:    api.API_VERSION,
			GoVersion:     runtime.Version(),
			BuildTime:     api.BUILDTIME,
			Os:            runtime.GOOS,
			Arch:          runtime.GOARCH,
			KernelVersion: kernelVersion,
			Experimental:  utils.ExperimentalBuild(),
		},
	}

	serverVersion, err := agentCli.Client().ServerVersion(ctx)
	if err == nil {
		vd.Server = &serverVersion
	}

	if err2 := tmpl.Execute(os.Stdout, vd); err2 != nil && err == nil {
		err = err2
	}
	fmt.Printf("\n")
	return err
}
