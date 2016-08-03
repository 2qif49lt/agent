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
	TlsKeyFile  string // tls连接所需
	TlsCertFile string // tls连接所需

	RsaSignKeyFile string // agentd 签名密钥 agent端有效
	RsaSignPubFile string // agent签名密钥  agentd端有效
}

var ComCfg *CommonFlags = nil

// InitCommonFlags initializes flags common to both client and daemon
func InitCommonFlags() *CommonFlags {
	var com = &CommonFlags{FlagSet: new(flag.FlagSet)}

	fs := com.FlagSet

	fs.BoolVarP(&com.Debug, "debug", "D", false, "Enable debug mode")
	fs.StringVarP(&com.LogLevel, "log-level", "L", "InfoLevel", "Set the logging level")

	fs.BoolVar(&com.NoTLS, "notls", false, "DO NOT Use TLS")
	fs.BoolVar(&com.NoSign, "nosign", false, "DO NOT verify the remote parameter sign")

	fs.StringVar(&tlsOptions.TlsCaFile, "tlsca", filepath.Join(certPath, DeafultTlsCaFile), "Trust certs signed only by this CA")

	fs.StringVar(&tlsOptions.TlsKeyFile, "tlskey", filepath.Join(certPath, DefaultTlsKeyFile), "Path to TLS key file")
	fs.StringVar(&tlsOptions.TlsCertFile, "tlscert", filepath.Join(certPath, DefaultCertFile), "Path to TLS certificate file")

	fs.StringVar(&tlsOptions.RsaSignKeyFile, "rsakey", filepath.Join(certPath, DefaultRsaSignPriFile), "Path to Rsa Sign private key file")
	fs.StringVar(&tlsOptions.RsaSignPubFile, "rsapub", filepath.Join(certPath, DefaultRsaSignPubFile), "Path to Rsa Sign public key file")

	fs.StringVarP(&tlsOptions.Host, "host", "H", fmt.Sprintf(`:%d`, DefaultAgentdListenPort), "Agentd listen address or Agent connect to,[ip]:port")

	/*
		-H参数可以为一下格式:
			1.master://target agent id:{port/name}. 表示要通过master 中转请求到目标服务器上的port服务,默认为agent.
			2.tcp://target agent ip:{port}.
			3.unix:///var/run/agentd.sock
			4.npipe:////./pipe/agentd_engine
			5.tcp://127.0.0.1:{port}
			6.:%d

		当-H参数未填充为默认值:%d.
			作为agentd:会默认监听所有地址.
			作为client:则会时则会自动选择3,4,5
		当-H参数不为默认值时,程序则根据实际参数进行监听或连接.
	*/
	ComCfg = com
	return com
}
