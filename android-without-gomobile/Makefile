all: build

build: clean
	CGO_ENABLED=1 GOOS=android GOARCH=arm64 CC=$$ANDROID_HOME/ndk-bundle/toolchains/llvm/prebuilt/linux-x86_64/bin/aarch64-linux-android21-clang CXX=$$ANDROID_HOME/ndk-bundle/toolchains/llvm/prebuilt/linux-x86_64/bin/aarch64-linux-android21-clang++ go build -buildmode=c-shared -o=out/arm64-v8a/libexample.so ./pkg/main/

clean:
	rm -rf out