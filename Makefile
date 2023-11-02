test:
	go test ./... -v

build:
	go build *.go

server: build
	./main

format:
	go fmt