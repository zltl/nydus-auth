
VERSION=0.0.1
GIT_REF=$(shell git describe --tags)
BUILD_TIME=$(shell date --rfc-3339=seconds | sed 's/ /T/')
GIT_HASH=$(shell git rev-parse HEAD)

LDFLAGS = -ldflags '-X version.Version=${VERSION} -X version.BuildTime=${BUILD_TIME} -X version.GitRef=${GIT_REF} -X version.GitHash=${GIT_HASH}'

.PHONY: build clean test

.ONESHELL:

build:
	go build $(LDFLAGS) ./cmd/...

test:
	go test ./...

clean:
	 rm -rf nydus-auth

