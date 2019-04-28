SHELL = /bin/bash
GIT_COMMIT := $(shell git rev-parse --short HEAD)
VERSION_FILE=VERSION.md
VERSION=$(shell cat ${VERSION_FILE})

GO_LDFLAGS := '-X "github.com/clintjedwards/cursor/cmd.appVersion=$(VERSION) $(GIT_COMMIT)"'

build:
	go build -ldflags $(GO_LDFLAGS) -o $(path)

run:
	go mod tidy
	go build -ldflags $(GO_LDFLAGS) -o /tmp/cursor && /tmp/cursor master

install:
	go build -ldflags $(GO_LDFLAGS) -o /tmp/cursor
	sudo mv /tmp/cursor /usr/local/bin/
	chmod +x /usr/local/bin/cursor
