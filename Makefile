testing_dependencies:
	# install golint
	go get -u github.com/golang/lint/golint

	# install gometalinter
	go get -u github.com/alecthomas/gometalinter

	# install all known linters:
	gometalinter --install

build:
	go build

lint:
	gometalinter --config="linters.json" ./...

test:
	go test ./...
