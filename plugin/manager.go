package plugin

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/2qif49lt/agent/api/types"
	"github.com/2qif49lt/agent/pkg/ioutils"
	"github.com/2qif49lt/logrus"
)

var manager *Manager

// ErrNotFound indicates that a plugin was not found locally.
type ErrNotFound string

func (name ErrNotFound) Error() string { return fmt.Sprintf("plugin %q not found", string(name)) }

// ErrInadequateCapability indicates that a plugin was found but did not have the requested capability.
type ErrInadequateCapability struct {
	name       string
	capability string
}

func (e ErrInadequateCapability) Error() string {
	return fmt.Sprintf("plugin %q found, but not with %q capability", e.name, e.capability)
}

type plugin struct {
	p                 types.Plugin
	stateSourcePath   string
	runtimeSourcePath string
}

func (p *plugin) Name() string {
	return p.p.Name
}

func (pm *Manager) newPlugin(name, id string) *plugin {
	p := &plugin{
		p: types.Plugin{
			Name: name,
			ID:   id,
		},
		stateSourcePath:   filepath.Join(pm.libRoot, id, "state"),
		runtimeSourcePath: filepath.Join(pm.runRoot, id),
	}
	return p
}

type pluginMap map[string]*plugin

// Manager controls the plugin subsystem.
type Manager struct {
	sync.RWMutex
	libRoot string
	runRoot string
	plugins pluginMap // TODO: figure out why save() doesn't json encode *plugin object

	nameToID map[string]string
}

// GetManager returns the singleton plugin Manager
func GetManager() *Manager {
	return manager
}

// Init (was NewManager) instantiates the singleton Manager.
// TODO: revert this to NewManager once we get rid of all the singletons.
func Init(root, execRoot string) (err error) {
	if manager != nil {
		return nil
	}

	root = filepath.Join(root, "plugins")
	execRoot = filepath.Join(execRoot, "plugins")
	for _, dir := range []string{root, execRoot} {
		if err := os.MkdirAll(dir, 0700); err != nil {
			return err
		}
	}

	manager = &Manager{
		libRoot:  root,
		runRoot:  execRoot,
		plugins:  make(map[string]*plugin),
		nameToID: make(map[string]string),
	}
	if err := os.MkdirAll(manager.runRoot, 0700); err != nil {
		return err
	}
	if err := manager.init(); err != nil {
		return err
	}
	return nil
}

func (pm *Manager) get(name string) (*plugin, error) {
	pm.RLock()
	id, nameOk := pm.nameToID[name]
	p, idOk := pm.plugins[id]
	pm.RUnlock()
	if !nameOk || !idOk {
		return nil, ErrNotFound(name)
	}
	return p, nil
}

func (pm *Manager) init() error {
	dt, err := os.Open(filepath.Join(pm.libRoot, "plugins.json"))
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	// TODO: Populate pm.plugins
	if err := json.NewDecoder(dt).Decode(&pm.nameToID); err != nil {
		return err
	}
	// FIXME: validate, restore

	return nil
}

func (pm *Manager) initPlugin(p *plugin) error {
	dt, err := os.Open(filepath.Join(pm.libRoot, p.p.ID, "manifest.json"))
	if err != nil {
		return err
	}
	err = json.NewDecoder(dt).Decode(&p.p.Manifest)
	dt.Close()
	if err != nil {
		return err
	}

	p.p.Config.Mounts = make([]types.PluginMount, len(p.p.Manifest.Mounts))
	for i, mount := range p.p.Manifest.Mounts {
		p.p.Config.Mounts[i] = mount
	}
	p.p.Config.Env = make([]string, 0, len(p.p.Manifest.Env))
	for _, env := range p.p.Manifest.Env {
		if env.Value != nil {
			p.p.Config.Env = append(p.p.Config.Env, fmt.Sprintf("%s=%s", env.Name, *env.Value))
		}
	}
	copy(p.p.Config.Args, p.p.Manifest.Args.Value)

	f, err := os.Create(filepath.Join(pm.libRoot, p.p.ID, "plugin-config.json"))
	if err != nil {
		return err
	}
	err = json.NewEncoder(f).Encode(&p.p.Config)
	f.Close()
	return err
}

func (pm *Manager) remove(p *plugin) error {
	if p.p.Active {
		return fmt.Errorf("plugin %s is active", p.p.Name)
	}
	pm.Lock() // fixme: lock single record
	defer pm.Unlock()
	os.RemoveAll(p.stateSourcePath)
	delete(pm.plugins, p.p.Name)
	pm.save()
	return nil
}
func (pm *Manager) enable(p *plugin) error {
	return fmt.Errorf("Not implemented")
}

func (pm *Manager) disable(p *plugin) error {
	return fmt.Errorf("Not implemented")
}

// fixme: not safe
func (pm *Manager) save() error {
	filePath := filepath.Join(pm.libRoot, "plugins.json")

	jsonData, err := json.Marshal(pm.nameToID)
	if err != nil {
		logrus.Debugf("Error in json.Marshal: %v", err)
		return err
	}
	ioutils.AtomicWriteFile(filePath, jsonData, 0600)
	return nil
}
