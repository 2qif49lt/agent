package events

import (
	"github.com/2qif49lt/agent/api/types/events"
	"github.com/2qif49lt/agent/api/types/filters"
)

// Filter can filter out docker events from a stream
type Filter struct {
	filter filters.Args
}

// NewFilter creates a new Filter
func NewFilter(filter filters.Args) *Filter {
	return &Filter{filter: filter}
}

// Include returns true when the event ev is included by the filters
func (ef *Filter) Include(ev events.Message) bool {
	return ef.filter.ExactMatch("event", ev.Action) &&
		ef.filter.ExactMatch("type", ev.Type) &&
		ef.matchDaemon(ev)
}

func (ef *Filter) matchDaemon(ev events.Message) bool {
	return ef.fuzzyMatchName(ev, events.DaemonEventType)
}

func (ef *Filter) fuzzyMatchName(ev events.Message, eventType string) bool {
	return ef.filter.FuzzyMatch(eventType, ev.Actor.ID) ||
		ef.filter.FuzzyMatch(eventType, ev.Actor.Attributes["name"])
}
