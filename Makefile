build:
	go build -o bin/goenvsubst main.go

run:
	go run main.go

fmt:
	go fmt ./...

test:
	go test ./...

all: build
