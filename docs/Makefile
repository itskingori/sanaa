SHELL := /bin/bash

install:
	bundle install

build:
	bundle exec jekyll build --baseurl=""

server:
	bundle exec jekyll serve --watch --safe --trace

test: build
	bundle exec rake test
