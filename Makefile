APP_NAME=warpqueue

.PHONY: help run build test lint fmt tidy clean

help:
	@echo "Available commands:"
	@echo " make run    - Run app"
	@echo " make build  - Build binary"
	@echo " make test   - Run tests"
	@echo " make lint   - Run lint checks"
	@echo " make fmt    - Format code"
	@echo " make tidy   - Clean dependencies"
	@echo " make clean  - Remove binary"

run:
	go run ./cmd

build:
	go build -o $(APP_NAME) ./cmd

test:
	go test ./... -v

lint:
	go vet ./...

fmt:
	go fmt ./...

tidy:
	go mod tidy

clean:
	rm -f $(APP_NAME)