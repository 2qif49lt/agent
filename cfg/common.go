package cfg

import (
	"fmt"
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
	Host  string // agentd:为监听地址. agent:请求的http地址,默认: ":1688"

	LogLevel string

	NoTLS  bool // 如为true,agentd不会建立tls服务.
	NoSign bool // 如为true,则agentd不会检查签名,agent不会携带签名.

	TlsCaFile   string // tls ca
	TlsKeyFile  string // agentd tls连接所需
	TlsCertFile string // agentd tls连接所需

	RsaSignKeyFile string // agentd 签名公钥
	RsaSignPubFile string // agent签名密钥
}

var ComCfg *CommonFlags = nil

// InitCommonFlags initializes flags common to both client and daemon
func InitCommonFlags() *CommonFlags {
	var com = &CommonFlags{FlagSet: new(flag.FlagSet)}

	fs := com.FlagSet

	fs.BoolVarP(&com.Debug, "debug", "D", false, "Enable debug mode")
	fs.StringVarP(&com.LogLevel, "log-level", "l", "InfoLevel", "Set the logging level")

	fs.BoolVar(&com.NoTLS, "notls", false, "DO NOT Use TLS")
	fs.BoolVar(&com.NoSign, "nosign", false, "DO NOT verify the remote parameter sign")

	fs.StringVar(&tlsOptions.TlsCaFile, "tlsca", filepath.Join(certPath, DeafultTlsCaFile), "Trust certs signed only by this CA")

	fs.StringVar(&tlsOptions.TlsKeyFile, "tlskey", filepath.Join(certPath, DefaultTlsKeyFile), "Path to TLS key file")
	fs.StringVar(&tlsOptions.TlsCertFile, "tlscert", filepath.Join(certPath, DefaultCertFile), "Path to TLS certificate file")

	fs.StringVar(&tlsOptions.RsaSignKeyFile, "rsapri", filepath.Join(certPath, DefaultRsaSignPriFile), "Path to Rsa Sign private key file")
	fs.StringVar(&tlsOptions.RsaSignPubFile, "rsapub", filepath.Join(certPath, DefaultRsaSignPubFile), "Path to Rsa Sign public key file")

	fs.StringVarP(&tlsOptions.Host, "host", "H", fmt.Sprintf(`:%d`, DefaultAgentdListenPort), "Agentd listen address or Agent connect to,[ip]:port")

	ComCfg = com
	return com
}
