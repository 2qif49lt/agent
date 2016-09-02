package plugin

import (
	"context"
	"net/http"

	"github.com/2qif49lt/agent/api/server/httputils"
	//	"github.com/2qif49lt/agent/api/types"
)

func (pr *pluginRouter) enablePlugin(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error {
	return pr.backend.Enable(vars["name"])
}

func (pr *pluginRouter) disablePlugin(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error {
	return pr.backend.Disable(vars["name"])
}

func (pr *pluginRouter) removePlugin(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error {
	return pr.backend.Remove(vars["name"])
}

func (pr *pluginRouter) listPlugins(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error {
	l, err := pr.backend.List()
	if err != nil {
		return err
	}
	return httputils.WriteJSON(w, http.StatusOK, l)
}

func (pr *pluginRouter) inspectPlugin(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error {
	result, err := pr.backend.Inspect(vars["name"])
	if err != nil {
		return err
	}
	return httputils.WriteJSON(w, http.StatusOK, result)
}
