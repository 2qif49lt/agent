package plugin

// Backend for Plugin
type Backend interface {
	Disable(name string) error
	Enable(name string) error
	List() ([]enginetypes.Plugin, error)
	Inspect(name string) (enginetypes.Plugin, error)
	Remove(name string) error
}
