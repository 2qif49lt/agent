package api

// Version 当前API 版本
const (
	SRV_VERSION = "0.2.0"
	CLI_VERSION = "0.1.0"

	// Version of Current REST API，比如稳定版本
	DefaultVersion string = "0.1.0"
	// MinVersion represents Minimum REST API version supported
	MinVersion = "0.1.0"
)

var (
	BUILDTIME = "20060102150405"
)

// go build -ldflags "-X verison.BUILDTIME=`date +%Y%m%d%H%M%S`" -o agent
// http://stackoverflow.com/questions/11354518/golang-application-auto-build-versioning
