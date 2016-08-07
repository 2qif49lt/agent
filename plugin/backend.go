// +build experimental

package plugin

import (
	"github.com/2qif49lt/agent/api/types"
)

// Disable deactivates a plugin, which implies that they cannot be used by containers.
func (pm *Manager) Disable(name string) error {
	p, err := pm.get(name)
	if err != nil {
		return err
	}
	return pm.disable(p)
}

// Enable activates a plugin, which implies that they are ready to be used by containers.
func (pm *Manager) Enable(name string) error {
	p, err := pm.get(name)
	if err != nil {
		return err
	}
	return pm.enable(p)
}

// Inspect examines a plugin manifest
func (pm *Manager) Inspect(name string) (tp types.Plugin, err error) {
	p, err := pm.get(name)
	if err != nil {
		return tp, err
	}
	return p.p, nil
}

/*
// Pull pulls a plugin and enables it.
func (pm *Manager) Pull(name string, metaHeader http.Header, authConfig *types.AuthConfig) (types.PluginPrivileges, error) {

}
*/
// List displays the list of plugins and associated metadata.
func (pm *Manager) List() ([]types.Plugin, error) {
	out := make([]types.Plugin, 0, len(pm.plugins))
	for _, p := range pm.plugins {
		out = append(out, p.p)
	}
	return out, nil
}

// Remove deletes plugin's root directory.
func (pm *Manager) Remove(name string) error {
	p, err := pm.get(name)
	if err != nil {
		return err
	}
	return pm.remove(p)
}
