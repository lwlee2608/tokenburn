GO = $(shell which go 2>/dev/null)

APP             := tokenburn
VERSION         ?= v0.1.0
LDFLAGS         := -ldflags "-X main.AppVersion=$(VERSION)"

.PHONY: all build clean run test

all: clean build

clean:
	$(GO) clean -testcache
	$(RM) -rf bin/*
build:
	$(GO) build -o bin/$(APP) $(LDFLAGS) cmd/$(APP)/*.go
run:
	$(GO) run $(LDFLAGS) cmd/$(APP)/*.go
test:
	$(GO) test -v ./...
