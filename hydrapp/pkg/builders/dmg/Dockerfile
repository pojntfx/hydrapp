FROM crazymax/osxcross:latest-debian AS osxcross

FROM golang:bookworm

COPY --from=osxcross /osxcross /osxcross

RUN rm -rf /opt
RUN ln -s /osxcross/macports/pkgs/opt/ /opt

# Install native dependencies
RUN apt update
RUN apt install -y genisoimage make clang lld libc6-dev bzip2 curl npm tree jq zlib1g-dev cmake

RUN curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh -s -- -y

ENV PATH="/osxcross/bin:/root/.cargo/bin:$PATH"
ENV LD_LIBRARY_PATH="/osxcross/lib:$LD_LIBRARY_PATH"

ENV OSXCROSS_PKG_CONFIG_USE_NATIVE_VARIABLES=1
ENV MACOSX_DEPLOYMENT_TARGET=12.0
ENV OSXCROSS_MP_INC=1
ENV GOOS=darwin
ENV CGO_ENABLED=1

RUN yes 1 | osxcross-macports select-mirror || true

RUN cargo install apple-codesign

RUN git clone --branch only_what_core_needs https://github.com/fanquake/libdmg-hfsplus.git /tmp/libdmg-hfsplus
WORKDIR /tmp/libdmg-hfsplus
RUN cmake . -B build
RUN make -C build/dmg -j "$(nproc)"
RUN install build/dmg/dmg /usr/local/bin
RUN rm -rf /tmp/libdmg-hfsplus

RUN git config --global --add safe.directory '*'

COPY entrypoint.sh /entrypoint.sh

WORKDIR /work

CMD /entrypoint.sh
