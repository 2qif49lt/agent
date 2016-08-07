package system

import (
	"time"

	"github.com/2qif49lt/agent/api/types"
	"github.com/2qif49lt/agent/api/types/events"
	"github.com/2qif49lt/agent/api/types/filters"
)

// Backend is the methods that need to be implemented to provide
// system specific functionality.
type Backend interface {
	SystemInfo() (*types.Info, error)
	SystemVersion() types.Version
	SubscribeToEvents(since, until time.Time, ef filters.Args) ([]events.Message, chan interface{})
	UnsubscribeFromEvents(chan interface{})
}
