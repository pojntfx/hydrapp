DESTDIR ?=
PREFIX ?= /usr/local

all: build

build:
	go build -ldflags=-linkmode=external -o out/desktop-integrated-webserver-rpm

install: build
	install -m 0755 -d $(DESTDIR)$(PREFIX)/bin
	install -m 0755 out/desktop-integrated-webserver-rpm $(DESTDIR)$(PREFIX)/bin