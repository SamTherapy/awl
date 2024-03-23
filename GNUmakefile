# SPDX-License-Identifier: BSD-3-Clause
# GNU Makefile allowing for building on Windows (with GNU Make)

include template.mk

ifeq ($(OS),Windows_NT)
	EXE := $(PROG).exe
else
	EXE := $(PROG)
ifeq ($(shell uname), Darwin)
	INSTALLFLAGS :=
else
	INSTALLFLAGS := D
endif
endif

## install: installs awl
.PHONY: install
ifeq ($(OS),Windows_NT)
install:
	$(GO) install $(GOFLAGS) .
else
install: all
	install -$(INSTALLFLAGS)m755 $(PROG) $(DESTDIR)$(PREFIX)/$(BIN)/$(PROG)
	install -$(INSTALLFLAGS)m644 docs/$(PROG).1 $(DESTDIR)$(MAN)/man1/$(PROG).1
# completions need to go in one specific place :)
	install -$(INSTALLFLAGS)m644 completions/bash.bash $(DESTDIR)$(PREFIX)/$(SHARE)/bash-completion/completions/$(PROG)
	install -$(INSTALLFLAGS)m644 completions/fish.fish $(DESTDIR)$(PREFIX)/$(SHARE)/fish/vendor_completions.d/$(PROG).fish
	install -$(INSTALLFLAGS)m644 completions/zsh.zsh $(DESTDIR)$(PREFIX)/$(SHARE)/zsh/site-functions/_$(PROG)
endif
