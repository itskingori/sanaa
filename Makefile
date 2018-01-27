SHELL := /bin/bash

dependencies:
	# install golint
	go get -u github.com/golang/lint/golint

	# install gometalinter
	go get -u github.com/alecthomas/gometalinter

	# install all known linters:
	gometalinter --install

install:
	glide install --strip-vendor

build:
	go build

lint:
	gometalinter --config="linters.json" ./...

test:
	go test -v ./...
