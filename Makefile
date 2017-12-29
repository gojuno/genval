all:
	go build
	go generate
	go test ./examples/...
