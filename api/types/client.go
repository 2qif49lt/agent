package types

import (
	"bufio"
	"net"

	"github.com/2qif49lt/agent/api/types/filters"
)

// EventsOptions hold parameters to filter events with.
type EventsOptions struct {
	Since   string
	Until   string
	Filters filters.Args
}

// HijackedResponse holds connection information for a hijacked request.
type HijackedResponse struct {
	Conn   net.Conn
	Reader *bufio.Reader
}

// Close closes the hijacked connection and reader.
func (h *HijackedResponse) Close() {
	h.Conn.Close()
}

// CloseWriter is an interface that implements structs
// that close input streams to prevent from writing.
type CloseWriter interface {
	CloseWrite() error
}

// CloseWrite closes a readWriter for writing.
func (h *HijackedResponse) CloseWrite() error {
	if conn, ok := h.Conn.(CloseWriter); ok {
		return conn.CloseWrite()
	}
	return nil
}

// VersionResponse holds version information for the client and the server
type VersionResponse struct {
	Client *Version
	Server *Version
}

// ServerOK returns true when the client could connect to the docker server
// and parse the information received. It returns false otherwise.
func (v VersionResponse) ServerOK() bool {
	return v.Server != nil
}

// ServiceCreateResponse contains the information returned to a client
// on the  creation of a new service.
type ServiceCreateResponse struct {
	// ID is the ID of the created service.
	ID string
}

// ServiceListOptions holds parameters to list  services with.
type ServiceListOptions struct {
	Filter filters.Args
}

// TaskListOptions holds parameters to list  tasks with.
type TaskListOptions struct {
	Filter filters.Args
}
