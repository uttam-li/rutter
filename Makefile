BINARY  := rutter
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
COMMIT  ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo none)
DATE    ?= $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
LDFLAGS := -s -w \
	-X github.com/uttam-li/rutter/internal/version.Version=$(VERSION) \
	-X github.com/uttam-li/rutter/internal/version.Commit=$(COMMIT) \
	-X github.com/uttam-li/rutter/internal/version.Date=$(DATE)

.PHONY: build test lint fmt check check-release hooks snapshot clean

build:
	go build -trimpath -ldflags "$(LDFLAGS)" -o bin/$(BINARY) ./cmd/rutter

test:
	go test ./...

lint:
	golangci-lint run ./...

fmt:
	gofmt -w .

check:
	./scripts/check.sh

check-release:
	@command -v goreleaser >/dev/null 2>&1 || { echo "goreleaser v2 is required: https://goreleaser.com/install/" >&2; exit 1; }
	goreleaser check

hooks:
	git config core.hooksPath .githooks
	chmod +x .githooks/pre-push
	@echo "pre-push gate installed"

snapshot:
	goreleaser release --snapshot --clean

clean:
	rm -rf bin dist
