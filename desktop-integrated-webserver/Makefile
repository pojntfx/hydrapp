SIZES ?= 16x16 22x22 24x24 32x32 36x36 48x48 64x64 72x72 96x96 128x128 192x192 256x256 512x512 1024x1024

build: clean
	go build -o out/desktop-integrated-webserver
	for icon in $(SIZES); do convert icon.png -resize $${icon} out/icon-$${icon}.png; done

install:
	install -D -m 0755 out/desktop-integrated-webserver /usr/bin/desktop-integrated-webserver
	install -D -m 0755 desktop-integrated-webserver.desktop /usr/share/applications/desktop-integrated-webserver.desktop
	install -D -m 0644 desktop-integrated-webserver.metainfo.xml /usr/share/metainfo/desktop-integrated-webserver.metainfo.xml
	for icon in $(SIZES); do install -D -m 0644 out/icon-$${icon}.png /usr/share/icons/hicolor/$${icon}/apps/desktop-integrated-webserver.png; done

uninstall:
	rm /usr/bin/desktop-integrated-webserver
	rm /usr/share/applications/desktop-integrated-webserver.desktop
	rm /usr/share/metainfo/desktop-integrated-webserver.metainfo.xml
	for icon in $(SIZES); do rm /usr/share/icons/hicolor/$${icon}/apps/desktop-integrated-webserver.png; done

clean:
	rm -rf out
