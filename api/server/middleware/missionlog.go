package middleware

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"time"

	"github.com/2qif49lt/agent/api/server/httputils"
	"github.com/2qif49lt/agent/errors"
	"github.com/2qif49lt/agent/pkg/eventdb"
	"github.com/2qif49lt/agent/pkg/ioutils"
	"github.com/2qif49lt/agent/pkg/random"
	"github.com/2qif49lt/logrus"
)

// MissionMiddleware record  the request mission.return a response header "mid"
// THIS MIDDLEWARE SHOULD APPEND AT LAST
func MissionMiddleware(handler func(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error) func(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error {
		logrus.Debugln("MissionMiddleware enter")
		defer logrus.Debugln("MissionMiddleware leave")

		mid, err := random.GetGuid()
		if err != nil {
			logrus.Warnln("GetGuid fail", err)
			mid = fmt.Sprintf("ffffffff%s%010d", time.Now().Format("20060102150405"), rand.Intn(1e10))
		}
		ctx = context.WithValue(ctx, "mid", mid)
		w.Header().Set("mid", mid)

		command := httputils.CommandFromRequest(r)
		paras := r.RequestURI
		body := ""

		if r.Method == "POST" && httputils.CheckForText(r) == nil {
			maxBodySize := 4096 // 4KB
			if r.ContentLength <= int64(maxBodySize) {
				rbody := r.Body
				bufReader := bufio.NewReaderSize(rbody, maxBodySize)
				r.Body = ioutils.NewReadCloserWrapper(bufReader, func() error { return rbody.Close() })
				if b, e := bufReader.Peek(maxBodySize); e == io.EOF {
					body = string(b)
				}
			}
		}
		begtime := time.Now()

		err = handler(ctx, w, r, vars)

		cost := time.Since(begtime) / time.Millisecond

		version := httputils.VersionFromContext(ctx)
		ua := httputils.ValueFromContext(ctx, httputils.UAStringKey)

		fmt.Println(version, "faaaa", ua)
		if eventerr := eventdb.InsertMission(mid, version, command, paras, body, errors.Str(err), int(cost)); eventerr != nil {
			logrus.WithFields(logrus.Fields{
				"mid":     mid,
				"version": version,
				"command": command,
				"paras":   paras,
				"body":    body,
				"result":  errors.Str(err),
				"cost":    int(cost),
			}).Warnln(eventerr.Error())
		}

		return err
	}
}
