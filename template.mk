# SPDX-License-Identifier: BSD-3-Clause
# Template for the BSD/GNU makefiles

HASH ?= `git describe --tags --always --dirty --broken | sed 's/^v//;s/\([^-]*-g\)/r\1/;s/-/./g' || echo "UNKNOWN"`

SOURCES ?= $(shell find . -name "*.go" -type f ! -name '*_test*')
TEST_SOURCES ?= $(shell find . -name "*_test.go" -type f)

CGO_ENABLED ?= 0
GO ?= go
TEST ?= $(GO) test -race -cover
COVER ?= $(GO) tool cover
GOFLAGS ?= -trimpath -ldflags="-s -w -X=main.version=$(HASH)"
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
	$(SCDOC) <$? >$@

doc/wiki/$(PROG).1.md: doc/$(PROG).1
	pandoc --from man --to gfm -o $@ $?

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

coverage/coverage.out: $(TEST_SOURCES)
	$(TEST) -coverprofile=$@ ./...

.PHONY: test
## test: run go test
test: coverage/coverage.out

.PHONY: test-ci
test-ci:
	$(TEST) ./...

## fuzz: runs fuzz tests
fuzz: $(TEST_SOURCES)
	$(TEST) -fuzz=FuzzFlags -fuzztime 10000x ./cmd
	$(TEST) -fuzz=FuzzDig -fuzztime 10000x ./cmd
	$(TEST) -fuzz=FuzzParseArgs -fuzztime 10000x ./cmd

fuzz-ci: $(TEST_SOURCES)
	$(TEST) -fuzz=FuzzFlags -fuzztime 1000x ./cmd
	$(TEST) -fuzz=FuzzDig -fuzztime 1000x ./cmd
	$(TEST) -fuzz=FuzzParseArgs -fuzztime 1000x ./cmd

.PHONY: full_test
full_test: test fuzz

coverage/cover.html: coverage/coverage.out
	$(COVER) -func=$?
	$(COVER) -html=$? -o $@

## cover: generates test coverage, output as HTML
cover: coverage/cover.html

## clean: clean the build files
.PHONY: clean
clean:
	$(GO) clean
# Ignore errors if you remove something that doesn't exist
	rm -f doc/$(PROG).1
	rm -f coverage/cover*
	rm -rf vendor

## help: Prints this help message
.PHONY: help
help:
	@echo "Usage: "
	@sed -n 's/^##//p' $(MAKEFILE_LIST) | column -t -s ':' |  sed -e 's/^/ /'
