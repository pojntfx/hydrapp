build: clean
	go build -o out/deb-hello-world

install: build
	install -m 0755 out/deb-hello-world /usr/bin

clean:
	rm -rf out