# SPDX-License-Identifier: BSD-3-Clause
# GNU Makefile allowing for building on Windows (with GNU Make)

include template.mk

ifeq ($(OS),Windows_NT)
	EXE := $(PROG).exe
else
	EXE := $(PROG)
endif

## install: installs awl
.PHONY: install
ifeq ($(OS),Windows_NT)
install:
	$(GO) install $(GOFLAGS) .
else
install: all
	install -Dm755 $(PROG) $(DESTDIR)$(PREFIX)/$(BIN)/$(PROG)
	install -Dm644 doc/$(PROG).1 $(DESTDIR)$(MAN)/man1/$(PROG).1
# completions need to go in one specific place :)
	install -Dm644 completions/zsh.zsh $(DESTDIR)/$(PREFIX)/$(SHARE)/zsh/site-functions/_$(PROG)
endif
