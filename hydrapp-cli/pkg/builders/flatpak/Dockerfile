FROM golang:bullseye

# Install native dependencies
RUN apt update
RUN apt install -y flatpak-builder xz-utils

COPY entrypoint.sh /entrypoint.sh

WORKDIR /work

CMD /entrypoint.sh
