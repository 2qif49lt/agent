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

// InitCommonFlags initializes flags common to both client and daemon
func InitCommonFlags(daemon bool) *CommonFlags {
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

	if daemon {
		fs.StringVarP(&tlsOptions.Host, "host", "H", fmt.Sprintf(`:%d`, DefaultAgentdListenPort), "Agentd listen address,[ip]:port")
	} else {
		fs.StringVarP(&tlsOptions.Host, "host", "H", fmt.Sprintf(`127.0.0.1:%d`, DefaultAgentdListenPort), "Agent connect to,[ip]:port")
	}

	return com
}

// Merge is a helper function that merges n FlagSets into a single dest FlagSet
// In case of name collision between the flagsets it will apply
// the destination FlagSet's errorHandling behavior.
func Merge(dest *flag.FlagSet, flagsets ...*flag.FlagSet) error {
	for _, fset := range flagsets {
		if fset.formal == nil {
			continue
		}
		for k, f := range fset.formal {
			if _, ok := dest.formal[k]; ok {
				var err error
				if fset.name == "" {
					err = fmt.Errorf("flag redefined: %s", k)
				} else {
					err = fmt.Errorf("%s flag redefined: %s", fset.name, k)
				}
				fmt.Fprintln(fset.Out(), err.Error())
				// Happens only if flags are declared with identical names
				switch dest.errorHandling {
				case ContinueOnError:
					return err
				case ExitOnError:
					os.Exit(2)
				case PanicOnError:
					panic(err)
				}
			}
			newF := *f
			newF.Value = mergeVal{f.Value, k, fset}
			if dest.formal == nil {
				dest.formal = make(map[string]*Flag)
			}
			dest.formal[k] = &newF
		}
	}
	return nil
}
