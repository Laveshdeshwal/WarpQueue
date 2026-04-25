APP_NAME=warpqueue
SERVER_BIN=$(APP_NAME)-server
WORKER_BIN=$(APP_NAME)-worker

.PHONY: help run-server run-worker build build-server build-worker test test-verbose lint fmt tidy clean

help:
	@echo "Available commands:"
	@echo " make run-server    - Run the HTTP server"
	@echo " make run-worker    - Run the worker process"
	@echo " make build         - Build both binaries"
	@echo " make build-server  - Build the server binary"
	@echo " make build-worker  - Build the worker binary"
	@echo " make test          - Run all tests"
	@echo " make test-verbose  - Run all tests with verbose output"
	@echo " make lint          - Run lint checks"
	@echo " make fmt           - Format code"
	@echo " make tidy          - Clean dependencies"
	@echo " make clean         - Remove built binaries"

run-server:
	go run ./cmd/server

build:
	go build -o $(SERVER_BIN) ./cmd/server
	go build -o $(WORKER_BIN) ./cmd/worker

run-worker:
	go run ./cmd/worker

build-server:
	go build -o $(SERVER_BIN) ./cmd/server

build-worker:
	go build -o $(WORKER_BIN) ./cmd/worker

test:
	go test ./...

test-verbose:
	go test ./... -v

lint:
	go vet ./...

fmt:
	go fmt ./...

tidy:
	go mod tidy

clean:
	rm -f $(SERVER_BIN) $(WORKER_BIN)
