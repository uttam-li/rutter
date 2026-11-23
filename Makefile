BINARY  := rutter
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
COMMIT  ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo none)
DATE    ?= $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
LDFLAGS := -s -w \
	-X github.com/uttam-li/rutter/internal/version.Version=$(VERSION) \
	-X github.com/uttam-li/rutter/internal/version.Commit=$(COMMIT) \
	-X github.com/uttam-li/rutter/internal/version.Date=$(DATE)

.PHONY: build test lint fmt snapshot clean

build:
	go build -trimpath -ldflags "$(LDFLAGS)" -o bin/$(BINARY) ./cmd/rutter

test:
	go test ./...

lint:
	golangci-lint run ./...

fmt:
	gofmt -w .

snapshot:
	goreleaser release --snapshot --clean

clean:
	rm -rf bin dist
