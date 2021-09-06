DESTDIR ?=
PREFIX ?= /usr/local

SIZES ?= 16x16 22x22 24x24 32x32 36x36 48x48 64x64 72x72 96x96 128x128 192x192 256x256 512x512

all: build

build:
	go build -ldflags=-linkmode=external -o out/flatpak-integrated-webserver
	for icon in $(SIZES); do convert icon.png -resize $${icon} out/icon-$${icon}.png; done

install: build
	install -D -m 0755 out/flatpak-integrated-webserver $(DESTDIR)$(PREFIX)/bin/flatpak-integrated-webserver
	desktop-file-install --dir=$(DESTDIR)$(PREFIX)/share/applications com.pojtinger.flatpakIntegratedWebserver.desktop
	for icon in $(SIZES); do install -D -m 0644 out/icon-$${icon}.png $(DESTDIR)$(PREFIX)/share/icons/hicolor/$${icon}/apps/com.pojtinger.flatpakIntegratedWebserver.png; done

uninstall:
	rm $(DESTDIR)$(PREFIX)/bin/flatpak-integrated-webserver
	rm $(DESTDIR)$(PREFIX)/share/applications/com.pojtinger.flatpakIntegratedWebserver.desktop
	for icon in $(SIZES); do rm $(DESTDIR)$(PREFIX)/share/icons/hicolor/$${icon}/apps/com.pojtinger.flatpakIntegratedWebserver.png; done

run: build
	out/flatpak-integrated-webserver

clean:
	rm -rf out
