# SPDX-License-Identifier: BSD-3-Clause
# Template for the BSD/GNU makefiles

HASH ?= `git describe --always --dirty --broken | sed 's/^v//;s/\([^-]*-g\)/r\1/;s/-/./g' || echo "UNKNOWN"`

SOURCES ?= $(shell find . -name "*.go" -type f ! -name '*_test*')
TEST_SOURCES ?= $(shell find . -name "*_test.go" -type f)

CGO_ENABLED ?= 0
GO ?= go
TEST ?= $(GO) test -race -cover
COVER ?= $(GO) tool cover
GOFLAGS ?= -ldflags="-s -w -X=main.version=$(HASH)" -trimpath
DESTDIR :=

PREFIX ?= /usr/local
BIN ?= bin
SHARE ?= share

SCDOC ?= scdoc
MAN ?= $(PREFIX)/$(SHARE)/man

PROG ?= awl

# hehe
all: $(PROG) doc/$(PROG).1

$(PROG): $(SOURCES)
	$(GO) build -o $(EXE) $(GOFLAGS) .

doc/$(PROG).1: doc/$(PROG).1.scd
	$(SCDOC) <$< >$@

doc/wiki/$(PROG).1.md: doc/$(PROG).1
	pandoc --from man --to gfm -o $@ $<

## update_doc: update documentation (requires pandoc)
update_doc: doc/wiki/$(PROG).1.md

.PHONY: fmt
fmt:
	gofmt -w -s .

.PHONY: vet
vet:
	$(GO) vet ./...

## lint: lint awl, using fmt, vet and golangci-lint
.PHONY: lint
lint: fmt vet
	golangci-lint run --fix

## test: run go test
test: $(TEST_SOURCES)
	$(TEST) -v -coverprofile=coverage/coverage.out ./...

.PHONY: test-ci
test-ci:
	$(TEST) ./...

## fuzz: runs fuzz tests
fuzz: $(TEST_SOURCES)
	$(TEST) -fuzz=FuzzFlags -fuzztime 10000x ./cli
	$(TEST) -fuzz=FuzzDig -fuzztime 10000x ./cli
	$(TEST) -fuzz=FuzzParseArgs -fuzztime 10000x ./cli

fuzz-ci: $(TEST_SOURCES)
	$(TEST) -fuzz=FuzzFlags -fuzztime 1000x ./cli
	$(TEST) -fuzz=FuzzDig -fuzztime 1000x ./cli
	$(TEST) -fuzz=FuzzParseArgs -fuzztime 1000x ./cli

.PHONY: full_test
full_test: test fuzz

coverage/coverage.out: test
	$(COVER) -func=$@
	$(COVER) -html=$@ -o coverage/cover.html

## cover: generates test coverage, output as HTML
cover: coverage/coverage.out

## clean: clean the build files
.PHONY: clean
clean:
	$(GO) clean
# Ignore errors if you remove something that doesn't exist
	rm -f doc/$(PROG).1

## help: Prints this help message
.PHONY: help
help:
	@echo "Usage: "
	@sed -n 's/^##//p' $(MAKEFILE_LIST) | column -t -s ':' |  sed -e 's/^/ /'
