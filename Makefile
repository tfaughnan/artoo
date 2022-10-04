GO ?= go
GOFLAGS ?= -tags netgo

all: artoo

artoo:
	$(GO) build $(GOFLAGS) -o $@ cmd/artoo/main.go
clean:
	rm -f artoo

.PHONY: all artoo
