package cfg

import (
	"crypto/tls"
	"fmt"
	"github.com/2qif49lt/agent/pkg/connections/tlsconfig"
	flag "github.com/2qif49lt/pflag"
	"path/filepath"
)

/*
var ip = flag.IntP("flagname", "f", 1234, "help message")
flag.Lookup("flagname").NoOptDefVal = "4321"
Would result in something like

Parsed Arguments	Resulting Value
--flagname=1357	ip=1357
--flagname	ip=4321
[nothing]	ip=1234
Command line flag syntax

--flag    // boolean flags, or flags with no option default values
--flag x  // only on flags without a default value
--flag=x
*/
// CommonFlags are flags common to both the client and the daemon.
type CommonFlags struct {
	FlagSet *flag.FlagSet

	Debug bool   //
	Host  string // agentd:为监听地址. 如果地址不是127.0.0.1 则同时监听127.0.0.1

	LogLevel string

	NoTLS bool // 如为true,agentd不会建立tls服务.
	NoRsa bool // 如为true,则agentd不会检查签名,agent不会携带签名.

	TLSOptions *tlsconfig.Options

	RsaSignFile string // agentd 签名密钥 agent端有效 // agent签名密钥  agentd端有效

	Master string
}

var ComCfg *CommonFlags = nil

// InitCommonFlags initializes flags common to both client and daemon
func InitCommonFlags() *CommonFlags {
	if ComCfg != nil {
		return ComCfg
	}
	var com = &CommonFlags{FlagSet: new(flag.FlagSet)}

	fs := com.FlagSet

	fs.BoolVarP(&com.Debug, "debug", "D", false, "Enable debug mode")
	fs.StringVar(&com.LogLevel, "log-level", "", "Set the logging level")

	fs.BoolVar(&com.NoTLS, "notls", false, "DO NOT Use TLS")
	fs.BoolVar(&com.NoRsa, "norsa", false, "DO NOT verify the remote parameter sign")

	var tlsOptions tlsconfig.Options
	com.TLSOptions = &tlsOptions
	fs.StringVar(&com.TLSOptions.CAFile, "tlsca", filepath.Join(certPath, DeafultTlsCaFile), "Trust certs signed only by this CA")

	fs.StringVar(&com.TLSOptions.KeyFile, "tlskey", filepath.Join(certPath, DefaultTlsKeyFile), "Path to TLS key file")
	fs.StringVar(&com.TLSOptions.CertFile, "tlscert", filepath.Join(certPath, DefaultCertFile), "Path to TLS certificate file")
	//	fs.BoolVar(&com.TLSOptions.InsecureSkipVerify,"tlsskip", false, "controls whether a client verifies the server's certificate")

	fs.StringVar(&com.RsaSignFile, "rsakey", filepath.Join(certPath, DefaultRsaSignFile), "Path to Rsa Sign public/private key file")

	fs.StringVar(&tlsOptions.Host, "host", "", "Agentd listen address or Agent connect to,[ip]:port")

	fs.StringVar(&tlsOptions.Master, "master", "", "Address of master service")

	/*
		--host参数可以为一下格式:
			1.master://target agent id:{port/name}. 表示要通过master 中转请求到目标服务器上的port服务,默认为agent.
			2.tcp://target agent ip:{port}.
			3.unix:///var/run/agentd.sock
			4.npipe:////./pipe/agentd_engine
			5.tcp://127.0.0.1:{port}
			6.空

		当--host参数为空,作为agent时连127.0.0.1:3567
						作为agent时,且配置文件Host字段为空时:则监听127.0.0.1:3567

	*/
	ComCfg = com
	return com
}
func postParse(com *CommonFlags) {
	if com != nil {
		if com.NoTLS == true {
			com.TlsOptions = nil
		} else {
			com.TLSOptions.ClientAuth = tls.RequireAndVerifyClientCert
			com.TLSOptions.InsecureSkipVerify = false
		}
	}
}
