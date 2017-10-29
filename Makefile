build:
	go build

lint:
	gometalinter --config="linters.json" ./...

test:
	go test ./...
