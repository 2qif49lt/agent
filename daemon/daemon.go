// Package daemon exposes the functions that occur on the host server
// that the agent daemon is running.
//
// In implementing the various functions of the daemon, there is often
// a method-specific struct for configuring the runtime behavior.
package daemon

import (
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"sync"
	"syscall"
	"time"

	"github.com/2qif49lt/agent/api"
	//	"github.com/2qif49lt/agent/api/types"
	"github.com/2qif49lt/agent/daemon/events"
	"github.com/2qif49lt/agent/pkg/progress"
	"github.com/2qif49lt/agent/pkg/streamformatter"
	"github.com/2qif49lt/agent/pkg/system"
	"github.com/2qif49lt/agent/utils"
	"github.com/2qif49lt/logrus"

	"github.com/docker/libtrust"
)

var (
	errSystemNotSupported = fmt.Errorf("The Agent daemon is not supported on this platform.")
)

// Daemon holds information about the Docker daemon.
type Daemon struct {
	AgentID  string
	trustKey libtrust.PrivateKey

	configStore   *Config
	EventsService *events.Events
	root          string
	shutdown      bool
}

// 恢复状态
func (daemon *Daemon) restore() error {
	var (
		debug = utils.IsDebugEnabled()
	)

	if !debug {
		logrus.Info("restore: start.")
	}

	group := sync.WaitGroup{}
	for i := 0; i != 10; i++ {
		group.Add(1)

		go func(n int) {
			defer group.Done()

			logrus.Debugf("restore  %d", n)

		}(i)

	}
	group.Wait()

	for j := 0; j != 10; j++ {
		group.Add(1)
		go func(n int) {
			defer group.Done()
			logrus.Debugf("restore others  %d", n)

		}(j)
	}

	group.Wait()

	if !debug {
		if logrus.GetLevel() == logrus.InfoLevel {
			fmt.Println()
		}
		logrus.Info("restore: done.")
	}

	return nil
}

// NewDaemon sets up everything for the daemon to be able to service
// requests from the webserver.
func NewDaemon(config *Config) (daemon *Daemon, err error) {
	setDefaultMtu(config)

	// Ensure we have compatible and valid configuration options
	if err := verifyDaemonSettings(config); err != nil {
		return nil, err
	}

	// Verify the platform is supported as a daemon
	if !platformSupported {
		return nil, errSystemNotSupported
	}

	// Validate platform-specific requirements
	if err := checkSystem(); err != nil {
		return nil, err
	}

	// set up SIGUSR1 handler on Unix-like systems, or a Win32 global event
	// on Windows to dump Go routine stacks
	setupDumpStackTrap()

	d := &Daemon{configStore: config}
	// Ensure the daemon is properly shutdown if there is a failure during
	// initialization
	defer func() {
		if err != nil {
			if err := d.Shutdown(); err != nil {
				logrus.Error(err)
			}
		}
	}()

	if err := configureMaxThreads(config); err != nil {
		logrus.Warnf("Failed to configure golang's threads limit: %v", err)
	}

	trustKey, err := api.LoadOrCreateTrustKey(config.TrustKeyPath)
	if err != nil {
		return nil, err
	}

	trustDir := filepath.Join(config.Root, "trust")

	if err := system.MkdirAll(trustDir, 0700); err != nil {
		return nil, err
	}

	eventsService := events.New()

	//	d.AgentID = trustKey.PublicKey().KeyID()

	d.trustKey = trustKey
	d.idIndex = truncindex.NewTruncIndex([]string{})
	d.statsCollector = d.newStatsCollector(1 * time.Second)

	d.EventsService = eventsService

	d.root = config.Root

	if err := d.restore(); err != nil {
		return nil, err
	}

	return d, nil
}

// Shutdown stops the daemon.
func (daemon *Daemon) Shutdown() error {
	daemon.shutdown = true

	return nil
}

func writeDistributionProgress(cancelFunc func(), outStream io.Writer, progressChan <-chan progress.Progress) {
	progressOutput := streamformatter.NewJSONStreamFormatter().NewProgressOutput(outStream, false)
	operationCancelled := false

	for prog := range progressChan {
		if err := progressOutput.WriteProgress(prog); err != nil && !operationCancelled {
			// don't log broken pipe errors as this is the normal case when a client aborts
			if isBrokenPipe(err) {
				logrus.Info("Pull session cancelled")
			} else {
				logrus.Errorf("error writing progress to client: %v", err)
			}
			cancelFunc()
			operationCancelled = true
			// Don't return, because we need to continue draining
			// progressChan until it's closed to avoid a deadlock.
		}
	}
}

func isBrokenPipe(e error) bool {
	if netErr, ok := e.(*net.OpError); ok {
		e = netErr.Err
		if sysErr, ok := netErr.Err.(*os.SyscallError); ok {
			e = sysErr.Err
		}
	}
	return e == syscall.EPIPE
}

func setDefaultMtu(config *Config) {
	// do nothing if the config does not have the default 0 value.
	if config.Mtu != 0 {
		return
	}
	config.Mtu = defaultNetworkMtu
}

// IsShuttingDown tells whether the daemon is shutting down or not
func (daemon *Daemon) IsShuttingDown() bool {
	return daemon.shutdown
}

// Reload reads configuration changes and modifies the
// daemon according to those changes.
// These are the settings that Reload changes:
// - Daemon labels.
// - Daemon debug log level.
// - Daemon max concurrent downloads
// - Daemon max concurrent uploads
// - Cluster discovery (reconfigure and restart).
// - Daemon live restore
func (daemon *Daemon) Reload(config *Config) error {
	var err error
	// used to hold reloaded changes
	attributes := map[string]string{}

	// We need defer here to ensure the lock is released as
	// daemon.SystemInfo() will try to get it too
	defer func() {
		if err == nil {
			daemon.LogDaemonEventWithAttributes("reload", attributes)
		}
	}()

	return nil
}
