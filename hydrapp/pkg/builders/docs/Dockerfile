FROM debian:bookworm

# Install native dependencies
RUN apt update
RUN apt install -y pandoc

COPY entrypoint.sh /entrypoint.sh

WORKDIR /work

CMD /entrypoint.sh
