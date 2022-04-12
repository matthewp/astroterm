BINARY  := astroterm
TAG     := $(shell git describe --tags --always --abbrev=0 --match="v[0-9]*.[0-9]*.[0-9]*" 2> /dev/null)
VERSION := $(shell echo "${TAG}" | sed 's/^.//')

LDFLAGS := -ldflags "-X 'astroterm/version.Version=${VERSION}'"

.PHONY: build
build:
	#
	# ################################################################################
	# >>> TARGET: build
	# ################################################################################
	#
	go build ${LDFLAGS} -o build/${BINARY}