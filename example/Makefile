DESTDIR ?=
PREFIX ?= /usr/local
DST ?= out/com.pojtinger.felicitas.hydrapp.example

SIZES ?= 16x16 22x22 24x24 32x32 36x36 48x48 64x64 72x72 96x96 128x128 192x192 256x256 512x512

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
