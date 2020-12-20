GO ?= go
PROJECT := go-service-example
PLATFORMS=darwin linux
GOARCH := amd64

vendor:
	go mod tidy
	go mod vendor
.PHONY: vendor

linter:
	./bin/golangci-lint run ./...
.PHONY: linter

test: linter
	env CGO_ENABLED=1 go test -race ./...
.PHONY: test

build: test vendor
	go build -o build/${PROJECT} ./cmd/${PROJECT}
.PHONY: build

build_all:
	$(foreach GOOS, $(PLATFORMS), \
            $(shell \
            export GOOS=$(GOOS); \
            export GOARCH=$(GOARCH); \
            go build -o build/${PROJECT}-${GOOS}-${GOARCH} ./app/cmd; ))

run: build
	build/${PROJECT}
.PHONY: run

golangcilint.download:
	curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh
	go mod download
.PHONY: golangcilint.download

config.copy:
	cp configs/.env.dist ./.env
.PHONY: config.copy

setup: golangcilint.download config.copy
.PHONY: setup