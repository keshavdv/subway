.PHONY: default server client deps fmt clean all release-all assets client-assets server-assets contributors
export GOPATH:=$(shell pwd)

BUILDTAGS=debug
default: all

deps: assets
	go get -tags '$(BUILDTAGS)' -d -v subway/...

server: deps
	go install -v -tags '$(BUILDTAGS)' subway/main/subwayd

client: deps
	go install -v -tags '$(BUILDTAGS)' subway/main/subway

fmt:
	go fmt subway/...