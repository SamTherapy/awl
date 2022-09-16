# SPDX-License-Identifier: BSD-3-Clause
# BSD/POSIX makefile

include template.mk

EXE := $(PROG)

doc/$(PROG).1: doc/$(PROG).1.scd
	$(SCDOC) <doc/$(PROG).1.scd >$@

doc/wiki/$(PROG).1.md: doc/$(PROG).1
	pandoc --from man --to gfm -o $@ doc/$(PROG).1

## install: installs awl
.PHONY: install
install: all
	install -Dm755 $(PROG) $(DESTDIR)$(PREFIX)/$(BIN)/$(PROG)
	install -Dm644 doc/$(PROG).1 $(DESTDIR)$(MAN)/man1/$(PROG).1
# completions need to go in one specific place :)
	install -Dm644 completions/zsh.zsh $(DESTDIR)/$(PREFIX)/$(SHARE)/zsh/site-functions/_$(PROG)