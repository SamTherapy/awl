# SPDX-License-Identifier: BSD-3-Clause
# Template for the BSD/GNU makefiles

HASH ?= `git describe --always --dirty || echo "UNKNOWN"`
VER ?= "git-$(HASH)"

CGO_ENABLED ?= 0
GO ?= go
COVER ?= $(GO) tool cover
GOFLAGS ?= -ldflags "-s -w -X=main.version=$(VER)" -trimpath

PREFIX ?= /usr/local
BIN ?= bin

SCDOC ?= scdoc
MAN ?= $(PREFIX)/share/man

PROG ?= awl

# hehe
all: $(PROG) doc/$(PROG).1

doc/$(PROG).1: doc/wiki/$(PROG).1.md
	@cp doc/awl.1 doc/awl.bak
	$(SCDOC) <doc/wiki/$(PROG).1.md >doc/$(PROG).1 && rm doc/awl.bak || mv doc/awl.bak doc/awl.1


## test: run go test
test:
	$(GO) test -cover -coverprofile=coverage/coverage.out ./...

coverage/coverage.out: test
	$(COVER) -func=coverage/coverage.out
	$(COVER) -html=coverage/coverage.out -o coverage/cover.html

## cover: generates test coverage, output as HTML
cover: coverage/coverage.out

fmt:
	gofmt -w -s .

vet:
	$(GO) vet ./...

## lint: lint awl, using fmt, vet and golangci-lint
lint: fmt vet
	-golangci-lint run --fix

## clean: clean the build files
clean:
	$(GO) clean

## help: Prints this help message
help:
	@echo "Usage: "
	@sed -n 's/^##//p' $(MAKEFILE_LIST) | column -t -s ':' |  sed -e 's/^/ /'

.PHONY: clean lint test fmt vet help