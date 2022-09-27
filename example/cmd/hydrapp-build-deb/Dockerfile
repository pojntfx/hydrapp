FROM debian:bullseye

# Install native dependencies
RUN apt update
RUN apt install -y dpkg-dev apt-utils ubuntu-keyring debian-ports-archive-keyring pbuilder qemu-user-static

COPY .pbuilderrc /root/.pbuilderrc
COPY entrypoint.sh /entrypoint.sh

WORKDIR /work

CMD /entrypoint.sh
