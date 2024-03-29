BINARY?=astroterm
TAG     := $(shell git describe --tags --always --abbrev=0 --match="[0-9]*.[0-9]*.[0-9]*" 2> /dev/null)
VERSION := $(shell echo "${TAG}" | sed 's/^//')

LDFLAGS := -ldflags "-X 'astroterm/version.Version=${VERSION}'"

.PHONY: build
build:
	go build ${LDFLAGS} -o build/${BINARY}
