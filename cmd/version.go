package main

var (
	Version   = "0.1.0"
	Buildtime = "20060102150405"
)

// go build -ldflags "-X main.Buildtime=`date +%Y%m%d%H%M%S`" -o agent
// http://stackoverflow.com/questions/11354518/golang-application-auto-build-versioning
