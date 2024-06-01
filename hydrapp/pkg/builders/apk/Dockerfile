FROM golang:bookworm

# Install native dependencies
RUN apt update
RUN apt install -y curl unzip default-jdk-headless make zip fdroidserver npm

# Install the Android SDK and NDK
ENV ANDROID_BUILD_TOOLS_VERSION='34.0.0'
ENV ANDROID_API_VERSION='34'

RUN curl -L -o '/tmp/commandlinetools.zip' 'https://dl.google.com/android/repository/commandlinetools-linux-11076708_latest.zip'
RUN rm -rf '/tmp/cmdline-tools'
RUN unzip -d '/tmp/' '/tmp/commandlinetools.zip'

RUN rm -rf '${HOME}/Android/Sdk'
RUN mkdir -p '${HOME}/Android/Sdk'
RUN yes | /tmp/cmdline-tools/bin/sdkmanager "build-tools;${ANDROID_BUILD_TOOLS_VERSION}" 'cmdline-tools;latest' 'platform-tools' "platforms;android-${ANDROID_API_VERSION}" 'ndk-bundle' --sdk_root="${HOME}/Android/Sdk"

ENV HOME='/root'
ENV ANDROID_HOME="${HOME}/Android/Sdk"
ENV ANDROID_SDK_ROOT="${HOME}/Android/Sdk"
ENV ANDROID_NDK_ROOT="${HOME}/Android/Sdk/ndk-bundle"
ENV ANDROID_NDK_VERSION='30'
ENV PATH="${PATH}:${ANDROID_HOME}/platform-tools"

# Install the legacy Android NDK, but only if we're on x86_64 since there are no builds available for other architectures
ENV ANDROID_LEGACY_API_VERSION='8'
ENV ANDROID_LEGACY_NDK_VERSION='16'

ENV ANDROID_LEGACY_HOME="${ANDROID_HOME}/ndk/android-ndk-r10e"

RUN if [ "$(uname -m)" = "x86_64" ]; then curl -L -o '/tmp/sdk-legacy.zip' 'https://dl.google.com/android/repository/android-ndk-r10e-linux-x86_64.zip'; fi
RUN if [ "$(uname -m)" = "x86_64" ]; then rm -rf "${ANDROID_LEGACY_HOME}"; fi
RUN if [ "$(uname -m)" = "x86_64" ]; then unzip -d "${ANDROID_HOME}/ndk" '/tmp/sdk-legacy.zip'; fi

RUN git config --global --add safe.directory '*'

COPY entrypoint.sh /entrypoint.sh

WORKDIR /hydrapp/work

CMD /entrypoint.sh
