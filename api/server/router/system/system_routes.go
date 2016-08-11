package system

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/2qif49lt/agent/api"
	"github.com/2qif49lt/agent/api/server/httputils"
	//"github.com/2qif49lt/agent/api/types"
	"github.com/2qif49lt/agent/api/types/events"
	"github.com/2qif49lt/agent/api/types/filters"
	timetypes "github.com/2qif49lt/agent/api/types/time"
	"github.com/2qif49lt/agent/errors"
	"github.com/2qif49lt/agent/pkg/ioutils"
	"github.com/2qif49lt/logrus"
	"golang.org/x/net/context"
)

func optionsHandler(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error {
	w.WriteHeader(http.StatusOK)
	return nil
}

func pingHandler(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error {
	if err := httputils.ParseForm(r); err != nil {
		return err
	}
	msg := r.Form.Get("msg")

	_, err := w.Write([]byte(fmt.Sprintf(`OK: %s`, msg)))
	return err
}

func (s *systemRouter) getInfo(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error {
	info, err := s.backend.SystemInfo()
	if err != nil {
		return err
	}

	return httputils.WriteJSON(w, http.StatusOK, info)
}

func (s *systemRouter) getVersion(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error {
	info := s.backend.SystemVersion()
	info.APIVersion = api.DefaultVersion

	return httputils.WriteJSON(w, http.StatusOK, info)
}

func (s *systemRouter) getEvents(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error {
	if err := httputils.ParseForm(r); err != nil {
		return err
	}

	since, err := eventTime(r.Form.Get("since"))
	if err != nil {
		return err
	}
	until, err := eventTime(r.Form.Get("until"))
	if err != nil {
		return err
	}

	var (
		timeout        <-chan time.Time
		onlyPastEvents bool
	)
	if !until.IsZero() {
		if until.Before(since) {
			return errors.NewBadRequestError(fmt.Errorf("`since` time (%s) cannot be after `until` time (%s)", r.Form.Get("since"), r.Form.Get("until")))
		}

		now := time.Now()

		onlyPastEvents = until.Before(now)

		if !onlyPastEvents {
			dur := until.Sub(now)
			timeout = time.NewTimer(dur).C
		}
	}

	ef, err := filters.FromParam(r.Form.Get("filters"))
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	output := ioutils.NewWriteFlusher(w)
	defer output.Close()
	output.Flush()

	enc := json.NewEncoder(output)

	buffered, l := s.backend.SubscribeToEvents(since, until, ef)
	defer s.backend.UnsubscribeFromEvents(l)

	for _, ev := range buffered {
		if err := enc.Encode(ev); err != nil {
			return err
		}
	}

	if onlyPastEvents {
		return nil
	}

	for {
		select {
		case ev := <-l:
			jev, ok := ev.(events.Message)
			if !ok {
				logrus.Warnf("unexpected event message: %q", ev)
				continue
			}
			if err := enc.Encode(jev); err != nil {
				return err
			}
		case <-timeout:
			return nil
		case <-ctx.Done():
			logrus.Debug("Client context cancelled, stop sending events")
			return nil
		}
	}
}

func eventTime(formTime string) (time.Time, error) {
	t, tNano, err := timetypes.ParseTimestamps(formTime, -1)
	if err != nil {
		return time.Time{}, err
	}
	if t == -1 {
		return time.Time{}, nil
	}
	return time.Unix(t, tNano), nil
}
