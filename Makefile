# Project: gogitopsdeployer
# Makefile for standardizing build and development lifecycle

.PHONY: all build run test lint tidy clean install

# Binary name
APP_NAME=gogitopsdeployer
BINARY_PATH=bin/$(APP_NAME)

# Default target
all: tidy lint test build

build:
	@echo "Building binary..."
	@go build -o $(BINARY_PATH) cmd/agent/main.go

run:
	@echo "Running agent..."
	@go run cmd/agent/main.go

test:
	@echo "Running tests..."
	@go test -v ./...

lint:
	@echo "Running linter..."
	@go vet ./...

tidy:
	@echo "Tidying go modules..."
	@go mod tidy

clean:
	@echo "Cleaning binaries..."
	@rm -rf bin/

install:
	@echo "Installing binary..."
	@go install ./cmd/agent
