.PHONY: build install fmt vet test check clean

PREFIX ?= /usr/local
BINDIR = $(PREFIX)/bin
GO ?= $(shell which go)

# Attempt to find golangci-lint
LINTER := $(shell which golangci-lint 2> /dev/null || which $(HOME)/go/bin/golangci-lint 2> /dev/null || echo "golangci-lint")

build:
	$(GO) build -o bin/straw ./cmd/straw
	$(GO) build -o bin/strawd ./cmd/strawd

install: build
	install -d $(DESTDIR)$(BINDIR)
	install -m 755 bin/straw $(DESTDIR)$(BINDIR)/
	install -m 755 bin/strawd $(DESTDIR)$(BINDIR)/

fmt:
	$(GO) fmt ./...

vet:
	$(GO) vet ./...

test:
	$(GO) test ./...

check:
	$(GO) fmt ./...
	$(GO) vet ./...
	$(GO) test ./...
	$(LINTER) run ./...

clean:
	rm -rf bin/
	rm -f straw.log
