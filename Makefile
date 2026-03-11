.PHONY: build build-all test clean install docker docker-compose

APP_NAME=memorycli
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS=-ldflags "-s -w -X main.Version=$(VERSION)"

build:
	go build $(LDFLAGS) -o bin/$(APP_NAME).exe .

build-linux:
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o bin/$(APP_NAME)-linux-amd64 .

build-darwin:
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o bin/$(APP_NAME)-darwin-arm64 .

build-all:
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o bin/$(APP_NAME)-linux-amd64 .
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o bin/$(APP_NAME)-darwin-amd64 .
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o bin/$(APP_NAME)-darwin-arm64 .
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o bin/$(APP_NAME)-windows-amd64.exe .

test:
	go test -v ./...

clean:
	if exist bin rmdir /s /q bin

install: build
	copy bin\$(APP_NAME).exe C:\Windows\System32\

deps:
	go mod download
	go mod tidy

docker:
	docker build -t $(APP_NAME):$(VERSION) .

docker-compose:
	docker-compose up -d

redis:
	docker run -d --name memorycli-redis -p 6379:6379 redis/redis-stack-server:latest

redis-stop:
	docker stop memorycli-redis
	docker rm memorycli-redis

run: build
	.\bin\$(APP_NAME).exe

help:
	@echo "MemoryCLI Makefile"
	@echo ""
	@echo "Usage:"
	@echo "  make build        Build for current platform"
	@echo "  make build-all    Build for all platforms"
	@echo "  make test         Run tests"
	@echo "  make clean        Clean build artifacts"
	@echo "  make install      Install to system"
	@echo "  make redis        Start Redis container"
	@echo "  make docker       Build Docker image"
