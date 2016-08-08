package system

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
	"time"

	"golang.org/x/net/context"

	"github.com/2qif49lt/agent/api/types"
	eventtypes "github.com/2qif49lt/agent/api/types/events"
	"github.com/2qif49lt/agent/api/types/filters"
	"github.com/2qif49lt/agent/cli"
	"github.com/2qif49lt/agent/client"
	"github.com/2qif49lt/agent/pkg/jsonlog"
	"github.com/2qif49lt/cobra"
)

type eventsOptions struct {
	since  string
	until  string
	filter []string
}

// NewEventsCommand creats a new cobra.Command for `agent events`
func NewEventsCommand(agentCli *client.AgentCli) *cobra.Command {
	var opts eventsOptions

	cmd := &cobra.Command{
		Use:   "events [OPTIONS]",
		Short: "Get real time events from the server",
		Args:  cli.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runEvents(agentCli, &opts)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&opts.since, "since", "", "Show all events created since timestamp")
	flags.StringVar(&opts.until, "until", "", "Stream events until this timestamp")
	flags.StringSliceVarP(&opts.filter, "filter", "f", []string{}, "Filter output based on conditions provided")

	return cmd
}

func runEvents(agentCli *client.AgentCli, opts *eventsOptions) error {
	eventFilterArgs := filters.NewArgs()

	// Consolidate all filter flags, and sanity check them early.
	// They'll get process in the daemon/server.
	for _, f := range opts.filter {
		var err error
		eventFilterArgs, err = filters.ParseFlag(f, eventFilterArgs)
		if err != nil {
			return err
		}
	}

	options := types.EventsOptions{
		Since:   opts.since,
		Until:   opts.until,
		Filters: eventFilterArgs,
	}

	responseBody, err := agentCli.Client().Events(context.Background(), options)
	if err != nil {
		return err
	}
	defer responseBody.Close()

	return streamEvents(responseBody, os.Stdout)
}

// streamEvents decodes prints the incoming events in the provided output.
func streamEvents(input io.Reader, output io.Writer) error {
	return DecodeEvents(input, func(event eventtypes.Message, err error) error {
		if err != nil {
			return err
		}
		printOutput(event, output)
		return nil
	})
}

// 处理单个event 记录
type eventProcessor func(event eventtypes.Message, err error) error

// printOutput prints all types of event information.
// Each output includes the event type, actor id, name and action.
// Actor attributes are printed at the end if the actor has any.
func printOutput(event eventtypes.Message, output io.Writer) {
	if event.TimeNano != 0 {
		fmt.Fprintf(output, "%s ", time.Unix(0, event.TimeNano).Format(jsonlog.RFC3339NanoFixed))
	} else if event.Time != 0 {
		fmt.Fprintf(output, "%s ", time.Unix(event.Time, 0).Format(jsonlog.RFC3339NanoFixed))
	}

	fmt.Fprintf(output, "%s %s %s", event.Type, event.Action, event.Actor.ID)

	if len(event.Actor.Attributes) > 0 {
		var attrs []string
		var keys []string
		for k := range event.Actor.Attributes {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			v := event.Actor.Attributes[k]
			attrs = append(attrs, fmt.Sprintf("%s=%s", k, v))
		}
		fmt.Fprintf(output, " (%s)", strings.Join(attrs, ", "))
	}
	fmt.Fprint(output, "\n")
}
