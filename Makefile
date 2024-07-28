build:
	go build -o bin/go-envsubst main.go

run:
	go run main.go

fmt:
	go fmt ./...

test:
	go test ./...

mod:
	go mod tidy

all: mod fmt test build
