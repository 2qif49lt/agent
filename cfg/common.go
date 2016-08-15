package cfg

import (
	"crypto/tls"
	"fmt"
	"github.com/2qif49lt/agent/cfg/cfgfile"
	"github.com/2qif49lt/agent/pkg/connections/tlsconfig"
	flag "github.com/2qif49lt/pflag"
	"os"
	"path/filepath"
	"strings"
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
	fs.StringVar(&com.TLSOptions.CertFile, "tlscert", filepath.Join(certPath, DefultTlsCertFile), "Path to TLS certificate file")
	//	fs.BoolVar(&com.TLSOptions.InsecureSkipVerify,"tlsskip", false, "controls whether a client verifies the server's certificate")

	fs.StringVar(&com.RsaSignFile, "rsakey", filepath.Join(certPath, DefaultRsaSignFile), "Path to Rsa Sign public/private key file")

	fs.StringVar(&com.Host, "host", "", "Agentd listen address or Agent connect to,[ip]:port")

	fs.StringVar(&com.Master, "master", "", "Address of master service")

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

// PostCheck 在参数解析后执行检查合并参数
func PostCheck() error {
	if ComCfg == nil {
		return fmt.Errorf(`common flags is not been initialized!`)
	} else {
		if ComCfg.NoTLS == true {
			ComCfg.TLSOptions = nil
		} else {
			ComCfg.TLSOptions.ClientAuth = tls.RequireAndVerifyClientCert
			ComCfg.TLSOptions.InsecureSkipVerify = false
		}
	}

	if Conf == nil {
		return fmt.Errorf(`config is not been loaded!`)
	}

	mergeCommonConfig(ComCfg, Conf)

	if err := isTlsLegal(ComCfg.NoTLS, ComCfg.NoRsa); err != nil {
		return err
	}
	return nil
}

// 如果命令行参数为空,则以配置文件值替换,如果配置文件也为空,则这里设置值
func mergeCommonConfig(com *CommonFlags, f *cfgfile.ConfigFile) {

	fn := func(to *string, from, def string) {
		*to = strings.TrimSpace(*to)
		from = strings.TrimSpace(from)

		if *to == "" {
			*to = from
			if *to == "" {
				*to = def
			}
		}
	}
	fn(&com.LogLevel, f.Loglvl, "InfoLevel")
	fn(&com.Host, f.Host, fmt.Sprintf(`127.0.0.1:%d`, DefaultAgentdListenPort))
	fn(&com.Master, f.Master.Srvs, "127.0.0.1:3678")
}

// IsTlsLegal return whether agentd install properly
func isTlsLegal(notls, norsa bool) error {
	isexist := func(fn string) bool {
		_, err := os.Stat(filepath.Join(GetCertPath(), fn))
		return err == nil || os.IsExist(err)
	}

	if !notls {
		if isexist(DeafultTlsCaFile) && isexist(DefaultTlsKeyFile) &&
			isexist(DefultTlsCertFile) {
			return nil
		} else {
			return fmt.Errorf("cert files not exist")
		}
	}

	if !norsa {
		if isexist(DefaultRsaSignFile) {
			return nil
		} else {
			return fmt.Errorf("rsa files not exist")
		}
	}
	return nil
}
