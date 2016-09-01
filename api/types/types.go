package types

// Version contains response of Remote API:
// GET "/version"
type Version struct {
	APIVersion    string
	GoVersion     string
	Os            string
	Arch          string
	KernelVersion string
	Experimental  bool
	BuildTime     string
}

// Info contains response of Remote API:
// GET "/info"
type Info struct {
	ID                string
	Debug             bool
	NFd               int
	NGoroutines       int
	SystemTime        string
	NEventsListener   int
	KernelVersion     string
	OperatingSystem   string
	OSType            string
	Architecture      string
	NCPU              int
	MemTotal          int64
	ExperimentalBuild bool
	ServerVersion     string
	HTTPProxy         string
	HTTPSProxy        string
	NoProxy           string
	Name              string
}

// AuthResponse contains response of Remote API:
// POST "/auth"
type AuthResponse struct {
	// Status is the authentication status
	Status string `json:"Status"`

	// IdentityToken is an opaque token used for authenticating
	// a user after a successful login.
	IdentityToken string `json:"IdentityToken,omitempty"`
}

// Ping contains response of Remote API: /ping
type Pong struct {
	msg string
}
