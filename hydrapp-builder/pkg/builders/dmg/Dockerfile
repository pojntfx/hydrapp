FROM crazymax/osxcross:latest-debian AS osxcross

FROM golang:bullseye

COPY --from=osxcross /osxcross /osxcross

RUN rm -rf /opt
RUN ln -s /osxcross/macports/pkgs/opt/ /opt

# Install native dependencies
RUN apt update
RUN apt install -y genisoimage make imagemagick clang lld libc6-dev bzip2 curl npm tree jq

RUN curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh -s -- -y

ENV PATH="/osxcross/bin:/root/.cargo/bin:$PATH"
ENV LD_LIBRARY_PATH="/osxcross/lib:$LD_LIBRARY_PATH"

ENV OSXCROSS_PKG_CONFIG_USE_NATIVE_VARIABLES=1
ENV MACOSX_DEPLOYMENT_TARGET=11.3.0
ENV OSXCROSS_MP_INC=1
ENV GOOS=darwin
ENV CGO_ENABLED=1

RUN yes 1 | osxcross-macports select-mirror || true

RUN cargo install apple-codesign

COPY entrypoint.sh /entrypoint.sh

WORKDIR /work

CMD /entrypoint.sh
