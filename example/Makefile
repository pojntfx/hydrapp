DESTDIR ?=
PREFIX ?= /usr/local
DST ?= out/com.pojtinger.felicitas.hydrapp.example

all: build

build:
	go build -o $(DST)

install: build
	install -D -m 0755 $(DST) $(DESTDIR)$(PREFIX)/bin/com.pojtinger.felicitas.hydrapp.example

uninstall:
	rm $(DESTDIR)$(PREFIX)/bin/com.pojtinger.felicitas.hydrapp.example

run: build
	out/com.pojtinger.felicitas.hydrapp.example

clean:
	rm -rf out .flatpak-builder build-dir repo
