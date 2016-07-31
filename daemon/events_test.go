package daemon

import (
	"testing"
	"time"

	"github.com/docker/docker/daemon/events"
	eventtypes "github.com/docker/engine-api/types/events"
)

func validateTestAttributes(t *testing.T, l chan interface{}, expectedAttributesToTest map[string]string) {
	select {
	case ev := <-l:
		event, ok := ev.(eventtypes.Message)
		if !ok {
			t.Fatalf("Unexpected event message: %q", ev)
		}
		for key, expected := range expectedAttributesToTest {
			actual, ok := event.Actor.Attributes[key]
			if !ok || actual != expected {
				t.Fatalf("Expected value for key %s to be %s, but was %s (event:%v)", key, expected, actual, event)
			}
		}
	case <-time.After(10 * time.Second):
		t.Fatalf("LogEvent test timed out")
	}
}
