GO ?= go
GOFLAGS ?= -tags netgo
GOSRC != find . -type f -name '*.go'
GOSRC += go.mod go.sum
PREFIX ?= /usr/local
BINDIR ?= $(PREFIX)/bin
MANDIR ?= $(PREFIX)/share/man
DOCDIR ?= $(PREFIX)/share/doc

artoo: $(GOSRC)
	$(GO) build $(GOFLAGS) -o $@ cmd/artoo/main.go

install: artoo artoo.1 artoo.toml.5 artoo.toml.example
	install -d $(DESTDIR)$(BINDIR) $(DESTDIR)$(MANDIR)/man1 \
		$(DESTDIR)$(MANDIR)/man5 $(DESTDIR)$(DOCDIR)/artoo
	install -m755 artoo $(DESTDIR)$(BINDIR)/artoo
	install -m644 artoo.1 $(DESTDIR)$(MANDIR)/man1/artoo.1
	install -m644 artoo.toml.5 $(DESTDIR)$(MANDIR)/man5/artoo.toml.5
	install -m644 artoo.toml.example $(DESTDIR)$(DOCDIR)/artoo/artoo.toml.example

clean:
	rm -f artoo

.PHONY: all install clean
