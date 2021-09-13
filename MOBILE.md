# Mobile Support

On mobile, a (system) WebView is probably the best choice. For now, only Android support is planned, as supporting iOS would mean to participate in and to support a oppressive system of software dictatorship under the premise of "security" and "privacy".

Here, `gomobile` could be used.

## `gomobile` and Android SDK/NDK Installation

The following could also be wrapped into a helper tool or Docker image; on Android, there is no way to properly avoid CGo.

```shell
$ export ANDROID_BUILD_TOOLS_VERSION=31.0.0
$ export ANDROID_API_VERSION=21
$ curl -L -o /tmp/commandlinetools.zip https://dl.google.com/android/repository/commandlinetools-linux-7583922_latest.zip
$ unzip -d /tmp/ /tmp/commandlinetools.zip
$ mkdir -p ~/Android/Sdk
$ yes | /tmp/cmdline-tools/bin/sdkmanager "build-tools;${ANDROID_BUILD_TOOLS_VERSION}" "cmdline-tools;latest" "platform-tools" "platforms;android-${ANDROID_API_VERSION}" "ndk-bundle" --sdk_root=$HOME/Android/Sdk
$ echo 'export ANDROID_HOME=$HOME/Android/Sdk' >>~/.bashrc
$ echo 'export ANDROID_SDK_ROOT=$HOME/Android/Sdk' >>~/.bashrc
$ echo 'export ANDROID_NDK_ROOT=$HOME/Android/Sdk/ndk-bundle' >>~/.bashrc
$ source ~/.bashrc
$ go env -w GO111MODULE=auto
$ go get golang.org/x/mobile/cmd/gomobile
$ gomobile init
$ go get -d golang.org/x/mobile/example/basic
# Manually build and install an APK
$ gomobile build -v -target android -androidapi ${ANDROID_API_VERSION} -o /tmp/basic.apk golang.org/x/mobile/example/basic
$ adb install /tmp/basic.apk
# Build & install an APK using gomobile
$ gomobile install -v -target android -androidapi ${ANDROID_API_VERSION} golang.org/x/mobile/example/basic
```
