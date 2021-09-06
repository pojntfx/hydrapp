DESTDIR ?=
PREFIX ?= /usr/local

all: build

build:
	go build -o out/flatpak-integrated-webserver

install: build
	install -D -m 0755 out/flatpak-integrated-webserver $(DESTDIR)$(PREFIX)/bin/flatpak-integrated-webserver
	desktop-file-install --dir=$(DESTDIR)$(PREFIX)/share/applications com.pojtinger.flatpakIntegratedWebserver.desktop

clean:
	rm -rf out
