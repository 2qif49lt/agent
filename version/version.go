package version

// Version 当前API 版本
const (
	SRV_VERSION = "0.1.0"
	CLI_VERSION = "0.1.0"
)

var (
	BUILDTIME = "20060102150405"
)

// go build -ldflags "-X verison.BUILDTIME=`date +%Y%m%d%H%M%S`" -o agent
// http://stackoverflow.com/questions/11354518/golang-application-auto-build-versioning
