FROM golang:bullseye

# Install native dependencies
RUN apt update
RUN apt install -y genisoimage make imagemagick

COPY entrypoint.sh /entrypoint.sh

WORKDIR /work

CMD /entrypoint.sh
