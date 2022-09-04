# SPDX-License-Identifier: BSD-3-Clause
# BSD/POSIX makefile

include template.mk

EXE := $(PROG)

## install: installs awl
.PHONY: install
install: all
	install -Dm755 $(PROG) $(DESTDIR)$(PREFIX)/$(BIN)/$(PROG)
	install -Dm644 doc/$(PROG).1 $(DESTDIR)$(MAN)/man1/$(PROG).1
# completions need to go in one specific place :)
	install -Dm644 completions/zsh.zsh $(DESTDIR)/$(PREFIX)/$(SHARE)/zsh/site-functions/_$(PROG)