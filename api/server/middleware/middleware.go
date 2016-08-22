package middleware

import (
	"net/http"

	"golang.org/x/net/context"
)

// Middleware is an interface to allow the use of ordinary functions as server api filters.
// Any struct that has the appropriate signature can be registered as a middleware.
type Middleware interface {
	// WrapHandler 将中间件的操作封装返回一个新的函数
	WrapHandler(func(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error) func(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error
}
