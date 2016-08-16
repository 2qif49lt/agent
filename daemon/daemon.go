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
	"runtime"
	"runtime/debug"
	"sync"
	"syscall"

	"github.com/2qif49lt/agent/api"
	//	"github.com/2qif49lt/agent/api/types"
	"github.com/2qif49lt/agent/cfg"
	"github.com/2qif49lt/agent/daemon/events"
	"github.com/2qif49lt/agent/pkg/progress"
	"github.com/2qif49lt/agent/pkg/streamformatter"
	"github.com/2qif49lt/agent/utils"
	"github.com/2qif49lt/logrus"

	"github.com/docker/libtrust"
)

var (
	errSystemNotSupported = fmt.Errorf("The Agent daemon is not supported on this platform.")
)

// Daemon holds information about the Docker daemon.
type Daemon struct {
	pubKey        libtrust.PublicKey
	configStore   *Config
	EventsService *events.Events
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
		logrus.Info("restore: done.")
	}

	return nil
}

// NewDaemon sets up everything for the daemon to be able to service
// requests from the webserver.
func NewDaemon(config *Config) (daemon *Daemon, err error) {
	// Ensure we have compatible and valid configuration options
	if err := verifyDaemonSettings(config); err != nil {
		return nil, err
	}

	// Verify the platform is supported as a daemon
	if !platformSupported {
		return nil, errSystemNotSupported
	}
	if config.Root == true {
		// Validate platform-specific requirements
		if err := checkSystem(); err != nil {
			return nil, err
		}
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

	pubkey, err := api.LoadSignPubKey(filepath.Join(cfg.GetCertPath(),
		cfg.DefaultRsaSignFile))
	if err != nil {
		return nil, err
	}

	d.pubKey = pubkey

	eventsService := events.New()
	d.EventsService = eventsService

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

func writeDistributionProgress(cancelFunc func(), outStream io.Writer,
	progressChan <-chan progress.Progress) {

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

// IsShuttingDown tells whether the daemon is shutting down or not
func (daemon *Daemon) IsShuttingDown() bool {
	return daemon.shutdown
}

// Reload 读取并且应用配置
func (daemon *Daemon) Reload(config *Config) error {
	var err error
	// used to hold reloaded changes
	attributes := map[string]string{}

	defer func() {
		if err == nil {
			daemon.LogDaemonEventWithAttributes("reload", attributes)
		}
	}()

	return nil
}

// configureMaxThreads sets the Go runtime max threads threshold
func configureMaxThreads(config *Config) error {
	if config.OSMaxThreadNum > 0 {
		cpus := runtime.NumCPU()
		maxThreads := cpus * config.OSMaxThreadNum
		debug.SetMaxThreads(maxThreads)

		logrus.Infof("Program threads limit set to %d", maxThreads)
	}
	return nil
}
