SHELL = /bin/bash
GIT_COMMIT := $(shell git rev-parse --short HEAD)
VERSION=0.1.0

GO_LDFLAGS := '-X "github.com/clintjedwards/cursor/cmd.appVersion=$(VERSION) $(GIT_COMMIT)"'

build-protos:
	protoc --go_out=plugins=grpc:. api/*.proto
	protoc --go_out=plugins=grpc:. plugin/proto/*.proto

build:
	go build -ldflags $(GO_LDFLAGS) -o $(path)

run:
	protoc --go_out=plugins=grpc:. api/*.proto
	protoc --go_out=plugins=grpc:. plugin/proto/*.proto
	go mod tidy
	go build -ldflags $(GO_LDFLAGS) -o /tmp/cursor && /tmp/cursor master

install:
	protoc --go_out=plugins=grpc:. api/*.proto
	protoc --go_out=plugins=grpc:. plugin/proto/*.proto
	go mod tidy
	go build -ldflags $(GO_LDFLAGS) -o /tmp/cursor
	sudo mv /tmp/cursor /usr/local/bin/
	chmod +x /usr/local/bin/cursor
