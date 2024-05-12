FROM debian:sid

# Install native dependencies
RUN apt update
RUN DEBIAN_FRONTEND=noninteractive apt install -y dpkg-dev apt-utils ubuntu-keyring debian-ports-archive-keyring pbuilder qemu-user-static

RUN git config --global --add safe.directory '*'

COPY .pbuilderrc /root/.pbuilderrc
COPY entrypoint.sh /entrypoint.sh

WORKDIR /hydrapp/work

CMD /entrypoint.sh
