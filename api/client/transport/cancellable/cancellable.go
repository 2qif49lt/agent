// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package cancellable provides helper function to cancel http requests.
package cancellable

import (
	"context"
	"io"
	"net/http"

	"github.com/2qif49lt/agent/api/client/transport"
)

func nop() {}

var (
	testHookContextDoneBeforeHeaders = nop
	testHookDoReturned               = nop
	testHookDidBodyClose             = nop
)

// Do sends an HTTP request with the provided transport.Sender and returns an HTTP response.
// If the client is nil, http.DefaultClient is used.
// If the context is canceled or times out, ctx.Err() will be returned.
func Do(client transport.Sender, req *http.Request) (*http.Response, error) {
	if client == nil {
		client = http.DefaultClient
	}

	var cancel context.CancelFunc
	ctx := req.Context()
	ctx, cancel = context.WithCancel(ctx)
	req = req.WithContext(ctx)

	type responseAndError struct {
		resp *http.Response
		err  error
	}
	result := make(chan responseAndError, 1)

	go func() {
		resp, err := client.Do(req)
		result <- responseAndError{resp, err}
	}()

	var resp *http.Response

	select {
	case <-ctx.Done():
		cancel()
		// Clean up after the goroutine calling client.Do:
		go func() {
			if r := <-result; r.resp != nil && r.resp.Body != nil {
				r.resp.Body.Close()
			}
		}()
		return nil, ctx.Err()
	case r := <-result:
		var err error
		resp, err = r.resp, r.err
		if err != nil {
			return resp, err
		}
	}

	c := make(chan struct{})
	go func() {
		select {
		case <-ctx.Done():
			cancel()
		case <-c:
			// The response's Body is closed.
			// 一些收听者模式下做些必要的清理
		}
	}()
	resp.Body = &notifyingReader{resp.Body, c}

	return resp, nil
}

// notifyingReader is an io.ReadCloser that closes the notify channel after
// Close is called or a Read fails on the underlying ReadCloser.
type notifyingReader struct {
	io.ReadCloser
	notify chan<- struct{}
}

func (r *notifyingReader) Read(p []byte) (int, error) {
	n, err := r.ReadCloser.Read(p)
	if err != nil && r.notify != nil {
		close(r.notify)
		r.notify = nil
	}
	return n, err
}

func (r *notifyingReader) Close() error {
	err := r.ReadCloser.Close()
	if r.notify != nil {
		close(r.notify)
		r.notify = nil
	}
	return err
}
