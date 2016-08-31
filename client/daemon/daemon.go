package daemon

import (
	"crypto/tls"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/2qif49lt/agent/api"
	apiserver "github.com/2qif49lt/agent/api/server"
	"github.com/2qif49lt/agent/api/server/middleware"
	"github.com/2qif49lt/agent/api/server/router"
	pluginrouter "github.com/2qif49lt/agent/api/server/router/plugin"
	systemrouter "github.com/2qif49lt/agent/api/server/router/system"
	"github.com/2qif49lt/agent/cfg"
	clientdaemon "github.com/2qif49lt/agent/daemon"
	"github.com/2qif49lt/agent/pkg/connections/tlsconfig"
	"github.com/2qif49lt/agent/pkg/listeners"
	"github.com/2qif49lt/agent/pkg/signal"
	"github.com/2qif49lt/agent/plugin"
	"github.com/2qif49lt/agent/utils"
	"github.com/2qif49lt/logrus"
	flag "github.com/2qif49lt/pflag"
)

// DaemonCli represents the daemon CLI.
type DaemonCli struct {
	*clientdaemon.Config
	*cfg.CommonFlags

	configFile *string
	api        *apiserver.Server
	d          *clientdaemon.Daemon
}

// NewDaemonCli returns a pre-configured daemon CLI
func NewDaemonCli() *DaemonCli {
	daemonConfig := new(clientdaemon.Config)

	return &DaemonCli{
		Config:      daemonConfig,
		CommonFlags: cfg.ComCfg,
	}
}
func (cli *DaemonCli) InitFlags(cmd *flag.FlagSet) {
	cli.Config.InstallFlags(cmd)
}

func (cli *DaemonCli) run() {
	stopc := make(chan bool)
	defer close(stopc)

	signal.Trap(func() {
		cli.stop()
		<-stopc // wait for daemonCli.run() to return
	})

	// The serve API routine never exits unless an error occurs
	// We need to start it as a goroutine and wait on it so
	// daemon doesn't exit
	serveAPIWait := make(chan error)
	go cli.api.Wait(serveAPIWait) // 开始服务

	// Daemon is fully initialized and handling API traffic
	// Wait for serve API to complete
	errAPI := <-serveAPIWait

	shutdownDaemon(cli.d, 15)

	if errAPI != nil {
		logrus.Errorf("Shutting down due to ServeAPI error: %v", errAPI)
	}
}
func (cli *DaemonCli) start() (err error) {
	loglev, err := logrus.ParseLevel(cli.CommonFlags.LogLevel)
	if err != nil {
		loglev = logrus.InfoLevel
	}
	logrus.SetLevel(loglev)

	if utils.ExperimentalBuild() {
		logrus.Warn("Running experimental build")
	}

	if err := setDefaultUmask(); err != nil {
		logrus.Errorf("Failed to set umask: %v", err)
		return err
	}

	serverConfig := &apiserver.Config{
		Logging:       true,
		SocketGroup:   cli.Config.SocketGroup,
		Version:       api.API_VERSION,
		CorsHeaders:   cli.Config.CorsHeaders,
		CertExtenAuth: cli.Config.CertExtenAuth,
	}

	if cli.CommonFlags.NoTLS == false {
		tlsOptions := cli.CommonFlags.TLSOptions
		if cli.Config.NoTLSClientVerify == true {
			tlsOptions.ClientAuth = tls.NoClientCert
		}
		tlsConfig, err := tlsconfig.Server(*tlsOptions)
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

		todo: 如果cli.CommonFlags.Host只指定了端口,则在LOOPBACK上为正常http服务器,并且没有参数校验

	*/
	apisrv := apiserver.New(serverConfig)
	cli.api = apisrv

	for i := 0; i < len(Hosts); i++ {
		var err error

		protoAddr := Hosts[i]
		protoAddrParts := strings.SplitN(protoAddr, "://", 2)
		if len(protoAddrParts) != 2 {
			err = fmt.Errorf("bad format %s, expected PROTO://ADDR", protoAddr)
			logrus.Errorf(err.Error())
			return err
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
		apisrv.Accept(protoAddrParts[1], ls)
	}

	if err := pluginInit(); err != nil {
		return err
	}

	d, err := clientdaemon.NewDaemon(cli.Config)
	if err != nil {
		logrus.Errorf("Error starting daemon: %v", err)
		return err
	}

	cli.initMiddlewares(apisrv, serverConfig)
	initRouter(apisrv, d)

	cli.d = d
	cli.setupConfigReloadTrap()

	logrus.Info("Daemon has completed initialization")

	logrus.WithFields(logrus.Fields{
		"version":   api.API_VERSION,
		"buildtime": api.BUILDTIME,
		"agentid":   cli.Config.AgentID,
	}).Info("Daemon start")

	return nil
}

func (cli *DaemonCli) reloadConfig() {
	reload := func(config *clientdaemon.Config) {
		if err := cli.d.Reload(config); err != nil {
			logrus.Errorf("Error reconfiguring the daemon: %v", err)
			return
		}
	}

	if err := clientdaemon.ReloadConfiguration(*cli.configFile, flag.CommandLine, reload); err != nil {
		logrus.Error(err)
	}
}

func (cli *DaemonCli) stop() {
	if cli.api != nil {
		cli.api.Close()
	}
}

// shutdownDaemon just wraps daemon.Shutdown() to handle a timeout in case
// d.Shutdown() is waiting too long to kill container or worst it's
// blocked there
func shutdownDaemon(d *clientdaemon.Daemon, timeout time.Duration) {
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

func loadDaemonCliConfig(config *clientdaemon.Config, flags *flag.FlagSet, commonConfig *cfg.CommonFlags, configFile string) (*clientdaemon.Config, error) {
	return config, nil
}

func initRouter(s *apiserver.Server, d *clientdaemon.Daemon) {
	routers := []router.Router{
		systemrouter.NewRouter(d),
		pluginrouter.NewRouter(plugin.GetManager()),
		// 路由
	}

	s.InitRouter(utils.IsDebugEnabled(), routers...)
}

func (cli *DaemonCli) initMiddlewares(s *apiserver.Server, cfg *apiserver.Config) {
	vm := middleware.NewVersionMiddleware(api.API_VERSION, api.DEF_VERSION, api.MIN_VERSION)
	s.UseMiddleware(vm)

	if cfg.CorsHeaders != "" {
		c := middleware.NewCORSMiddleware(cfg.CorsHeaders)
		s.UseMiddleware(c)
	}

	u := middleware.NewUserAgentMiddleware()
	s.UseMiddleware(u)
}

func pluginInit() error {
	procpath, err := utils.GetProcAbsDir()
	if err != nil {
		return err
	}
	return plugin.Init(filepath.Join(procpath, "plugin"), filepath.Join(procpath, "plugin"))
}
