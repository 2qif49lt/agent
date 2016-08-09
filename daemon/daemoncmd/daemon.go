package daemoncmd

import (
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/2qif49lt/agent/api"
	apiserver "github.com/2qif49lt/agent/api/server"
	"github.com/2qif49lt/agent/api/server/middleware"
	"github.com/2qif49lt/agent/api/server/router"
	pluginrouter "github.com/2qif49lt/agent/api/server/router/plugin"
	systemrouter "github.com/2qif49lt/agent/api/server/router/system"
	"github.com/2qif49lt/agent/cfg"
	"github.com/2qif49lt/agent/daemon"
	"github.com/2qif49lt/agent/pkg/connections/tlsconfig"
	"github.com/2qif49lt/agent/pkg/jsonlog"
	"github.com/2qif49lt/agent/pkg/listeners"
	"github.com/2qif49lt/agent/pkg/pidfile"
	"github.com/2qif49lt/agent/pkg/signal"
	"github.com/2qif49lt/agent/pkg/system"
	"github.com/2qif49lt/agent/plugin"
	"github.com/2qif49lt/agent/utils"
	"github.com/2qif49lt/agent/version"
	"github.com/2qif49lt/logrus"
	flag "github.com/2qif49lt/pflag"
)

// DaemonCli represents the daemon CLI.
type DaemonCli struct {
	*daemon.Config
	*cfg.CommonFlags

	configFile *string
	api        *apiserver.Server
	d          *daemon.Daemon
}

// NewDaemonCli returns a pre-configured daemon CLI
func NewDaemonCli() *DaemonCli {
	daemonConfig := new(daemon.Config)

	return &DaemonCli{
		Config:      daemonConfig,
		CommonFlags: cfg.ComCfg,
	}
}

func (cli *DaemonCli) start() (err error) {
	stopc := make(chan bool)
	defer close(stopc)

	// read config

	if cli.CommonFlags.Debug {
		utils.EnableDebug()
	}

	if utils.ExperimentalBuild() {
		logrus.Warn("Running experimental build")
	}

	if err := setDefaultUmask(); err != nil {
		return fmt.Errorf("Failed to set umask: %v", err)
	}

	if cli.Config.Pidfile != "" {
		pf, err := pidfile.New(cli.Config.Pidfile)
		if err != nil {
			return fmt.Errorf("Error starting daemon: %v", err)
		}
		defer func() {
			if err := pf.Remove(); err != nil {
				logrus.Error(err)
			}
		}()
	}

	serverConfig := &apiserver.Config{
		Logging:     true,
		SocketGroup: cli.Config.SocketGroup,
		Version:     verison.SRV_VERSION,
		CorsHeaders: cli.Config.CorsHeaders,
	}

	if cli.CommonFlags.NoTLS == false {
		tlsOptions := cli.CommonFlags.TLSOptions
		if cli.Config.NoTLSClientVerify == true {
			tlsOptions.ClientAuth = tls.NoClientCert
		}
		tlsConfig, err := tlsconfig.Server(tlsOptions)
		if err != nil {
			return err
		}
		serverConfig.TLSConfig = tlsConfig
	}
	if cli.CommonFlags.NoRsa == false {
		serverConfig.RSAVerify = cli.CommonFlags.RsaSignFile
	}

	Hosts := make([]string, 0)
	Hosts = append(Hosts, cli.CommonFlags.Host)
	/*
		if !strings.Contains(cli.CommonFlags.Host, "127.0.0.1") &&
			!strings.Contains(cli.CommonFlags.Host, "localhost") {
			_, port, err := net.SplitHostPort(cli.CommonFlags.Host)
			Hosts = append(Hosts, fmt.Sprintf(`127.0.0.1:%s`, port))
		}

	*/
	api := apiserver.New(serverConfig)
	cli.api = api

	for i := 0; i < len(Hosts); i++ {
		var err error

		protoAddr := Hosts[i]
		protoAddrParts := strings.SplitN(protoAddr, "://", 2)
		if len(protoAddrParts) != 2 {
			return fmt.Errorf("bad format %s, expected PROTO://ADDR", protoAddr)
		}

		proto := protoAddrParts[0]
		addr := protoAddrParts[1]

		// It's a bad idea to bind to TCP without tlsverify.
		if proto == "tcp" && (serverConfig.TLSConfig == nil || serverConfig.TLSConfig.ClientAuth != tls.RequireAndVerifyClientCert) {
			logrus.Warn("[!] DON'T BIND ON ANY IP ADDRESS WITHOUT enable TLS IF YOU DON'T KNOW WHAT YOU'RE DOING [!]")
		}

		ls, err := listeners.Init(proto, addr, serverConfig.SocketGroup, serverConfig.TLSConfig)
		if err != nil {
			return err
		}
		ls = wrapListeners(proto, ls)

		logrus.Debugf("Listener created for HTTP on %s (%s)", protoAddrParts[0], protoAddrParts[1])
		api.Accept(protoAddrParts[1], ls)
	}

	signal.Trap(func() {
		cli.stop()
		<-stopc // wait for daemonCli.start() to return
	})

	if err := pluginInit(); err != nil {
		return err
	}

	d, err := daemon.NewDaemon(cli.Config)
	if err != nil {
		return fmt.Errorf("Error starting daemon: %v", err)
	}

	logrus.Info("Daemon has completed initialization")

	logrus.WithFields(logrus.Fields{
		"version":   verison.SRV_VERSION,
		"buildtime": verison.BUILDTIME,
	}).Info("Agent daemon")

	cli.initMiddlewares(api, serverConfig)
	initRouter(api, d)

	cli.d = d
	cli.setupConfigReloadTrap()

	// The serve API routine never exits unless an error occurs
	// We need to start it as a goroutine and wait on it so
	// daemon doesn't exit
	serveAPIWait := make(chan error)
	go api.Wait(serveAPIWait) // 开始服务

	// after the daemon is done setting up we can notify systemd api
	notifySystem()

	// Daemon is fully initialized and handling API traffic
	// Wait for serve API to complete
	errAPI := <-serveAPIWait

	shutdownDaemon(d, 15)

	if errAPI != nil {
		return fmt.Errorf("Shutting down due to ServeAPI error: %v", errAPI)
	}

	return nil
}

func (cli *DaemonCli) reloadConfig() {
	reload := func(config *daemon.Config) {
		if err := cli.d.Reload(config); err != nil {
			logrus.Errorf("Error reconfiguring the daemon: %v", err)
			return
		}
	}

	if err := daemon.ReloadConfiguration(*cli.configFile, flag.CommandLine, reload); err != nil {
		logrus.Error(err)
	}
}

func (cli *DaemonCli) stop() {
	cli.api.Close()
}

// shutdownDaemon just wraps daemon.Shutdown() to handle a timeout in case
// d.Shutdown() is waiting too long to kill container or worst it's
// blocked there
func shutdownDaemon(d *daemon.Daemon, timeout time.Duration) {
	ch := make(chan struct{})
	go func() {
		d.Shutdown()
		close(ch)
	}()
	select {
	case <-ch:
		logrus.Debug("Clean shutdown succeeded")
	case <-time.After(timeout * time.Second):
		logrus.Error("Force shutdown daemon")
	}
}

func loadDaemonCliConfig(config *daemon.Config, flags *flag.FlagSet, commonConfig *cliflags.CommonFlags, configFile string) (*daemon.Config, error) {
	config.Debug = commonConfig.Debug
	config.Hosts = commonConfig.Hosts
	config.LogLevel = commonConfig.LogLevel
	config.TLS = commonConfig.TLS
	config.TLSVerify = commonConfig.TLSVerify
	config.CommonTLSOptions = daemon.CommonTLSOptions{}

	if commonConfig.TLSOptions != nil {
		config.CommonTLSOptions.CAFile = commonConfig.TLSOptions.CAFile
		config.CommonTLSOptions.CertFile = commonConfig.TLSOptions.CertFile
		config.CommonTLSOptions.KeyFile = commonConfig.TLSOptions.KeyFile
	}

	if configFile != "" {
		c, err := daemon.MergeDaemonConfigurations(config, flags, configFile)
		if err != nil {
			if flags.IsSet(daemonConfigFileFlag) || !os.IsNotExist(err) {
				return nil, fmt.Errorf("unable to configure the Docker daemon with file %s: %v\n", configFile, err)
			}
		}
		// the merged configuration can be nil if the config file didn't exist.
		// leave the current configuration as it is if when that happens.
		if c != nil {
			config = c
		}
	}

	if err := daemon.ValidateConfiguration(config); err != nil {
		return nil, err
	}

	// Regardless of whether the user sets it to true or false, if they
	// specify TLSVerify at all then we need to turn on TLS
	if config.IsValueSet(cliflags.TLSVerifyKey) {
		config.TLS = true
	}

	// ensure that the log level is the one set after merging configurations
	cliflags.SetDaemonLogLevel(config.LogLevel)

	return config, nil
}

func initRouter(s *apiserver.Server, d *daemon.Daemon) {
	routers := []router.Router{
		systemrouter.NewRouter(d),
		pluginrouter.NewRouter(d),
		// 路由
	}

	s.InitRouter(utils.IsDebugEnabled(), routers...)
}

func (cli *DaemonCli) initMiddlewares(s *apiserver.Server, cfg *apiserver.Config) {
	v := cfg.Version

	vm := middleware.NewVersionMiddleware(v, api.DefaultVersion, api.MinVersion)
	s.UseMiddleware(vm)

	if cfg.CorsHeaders != "" {
		c := middleware.NewCORSMiddleware(cfg.CorsHeaders)
		s.UseMiddleware(c)
	}

	u := middleware.NewUserAgentMiddleware(v)
	s.UseMiddleware(u)
}

func pluginInit() error {
	procpath, err := utils.GetProcAbsDir()
	if err != nil {
		return err
	}
	return plugin.Init(filepath.Join(procpath, "plugin"), filepath.Join(procpath, "plugin"))
}
