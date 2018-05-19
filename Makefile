SHELL := /bin/bash

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
	go build

lint:
	gometalinter --config="linters.json" ./...

test:
	go test -v ./...
