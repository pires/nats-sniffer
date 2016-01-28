GOPATH=$(shell pwd):$(shell pwd)/vendor

SRC=src/github.com/pires/nats-sniffer/
PKG=github.com/pires/nats-sniffer/...

.DEFAULT_GOAL := build

all: test build

build: clean
	@GOARCH=amd64 gb build all

.PHONY: clean
clean:
	@rm -rf ./{bin,pkg}

.PHONY: linux
linux: test
	@GOOS=linux GOARCH=amd64 gb build -ldflags '-w -extldflags=-static'

.PHONY: test
test:
	@go test -v ${PKG}
