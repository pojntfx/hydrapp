DESTDIR ?=
PREFIX ?= /usr/local
DST ?= out/com.pojtinger.felix.hydrapp.example

all: build

build:
	go build -o $(DST)

install: build
	install -D -m 0755 $(DST) $(DESTDIR)$(PREFIX)/bin/com.pojtinger.felix.hydrapp.example

uninstall:
	rm $(DESTDIR)$(PREFIX)/bin/com.pojtinger.felix.hydrapp.example

run: build
	out/com.pojtinger.felix.hydrapp.example

clean:
	rm -rf out .flatpak-builder build-dir repo
