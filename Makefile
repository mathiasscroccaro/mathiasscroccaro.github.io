test:
	go test ./... -v

build:
	go build *.go

server: build
	./main -s

format:
	go fmt