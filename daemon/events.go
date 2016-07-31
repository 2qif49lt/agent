package daemon

import (
	"strings"
	"time"

	"github.com/2qif49lt/agent/api/types/events"
	"github.com/2qif49lt/agent/api/types/filters"
	daemonevents "github.com/2qif49lt/agent/daemon/events"
)

// LogDaemonEventWithAttributes generates an event related to the daemon itself with specific given attributes.
func (daemon *Daemon) LogDaemonEventWithAttributes(action string, attributes map[string]string) {
	if daemon.EventsService != nil {
		if info, err := daemon.SystemInfo(); err == nil && info.Name != "" {
			attributes["name"] = info.Name
		}
		actor := events.Actor{
			ID:         daemon.AgentID,
			Attributes: attributes,
		}
		daemon.EventsService.Log(action, events.DaemonEventType, actor)
	}
}

// SubscribeToEvents returns the currently record of events, a channel to stream new events from, and a function to cancel the stream of events.
func (daemon *Daemon) SubscribeToEvents(since, until time.Time, filter filters.Args) ([]events.Message, chan interface{}) {
	ef := daemonevents.NewFilter(filter)
	return daemon.EventsService.SubscribeTopic(since, until, ef)
}

// UnsubscribeFromEvents stops the event subscription for a client by closing the
// channel where the daemon sends events to.
func (daemon *Daemon) UnsubscribeFromEvents(listener chan interface{}) {
	daemon.EventsService.Evict(listener)
}
