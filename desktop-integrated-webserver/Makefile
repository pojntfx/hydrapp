build: clean
	go build -o out/desktop-integrated-webserver

install: build
	install -m 0755 out/desktop-integrated-webserver /usr/bin

clean:
	rm -rf out