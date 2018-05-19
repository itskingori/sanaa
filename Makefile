SHELL := /bin/bash

VERSION := $(version)
MAJOR_VERSION := $(shell echo $(VERSION) | cut -d'.' -f1)
MINOR_VERSION := $(shell echo $(VERSION) | cut -d'.' -f2)
PATCH_VERSION := $(shell echo $(VERSION) | cut -d'.' -f3)
COMMIT_VERSION := $(shell git rev-parse HEAD)
PACKAGE_PATH := github.com/itskingori/sanaa
LDFLAGS := \
-X $(PACKAGE_PATH)/service.majorVersion=$(MAJOR_VERSION) \
-X $(PACKAGE_PATH)/service.minorVersion=$(MINOR_VERSION) \
-X $(PACKAGE_PATH)/service.patchVersion=$(PATCH_VERSION) \
-X $(PACKAGE_PATH)/service.commitVersion=$(COMMIT_VERSION)

.PHONY: all tools dependencies build lint test

all: dependencies build

tools:
	# install golint
	go get -u github.com/golang/lint/golint

	# install gometalinter
	go get -u github.com/alecthomas/gometalinter

	# install gox
	go get -v github.com/mitchellh/gox

	# install all known linters:
	gometalinter --install

dependencies:
	dep ensure

build:
	go build -ldflags="$(LDFLAGS)"

lint:
	gometalinter --config="linters.json" ./...

test:
	go test -v ./...
