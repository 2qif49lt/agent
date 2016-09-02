.PHONY: build binary

APP_NAME := agent
APP_FOLDER := binary

BUILD_TIME :=`date +%Y%m%d%H%M%S`
LDFLAGS = -ldflags "-X github.com/2qif49lt/agent/api.BUILDTIME=$(BUILD_TIME)"

build: binary
	go build $(LDFLAGS)  -o ./$(APP_FOLDER)/$(APP_NAME) ./main/*.go
	@echo "done"
binary:  
	mkdir -p binary
	
