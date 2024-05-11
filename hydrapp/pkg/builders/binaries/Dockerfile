FROM golang:bookworm

# Install native dependencies
RUN apt update
RUN DEBIAN_FRONTEND=noninteractive apt install -y curl npm tree jq

RUN curl -L -o /tmp/bagop "https://github.com/pojntfx/bagop/releases/latest/download/bagop.linux-$(uname -m)"
RUN install /tmp/bagop /usr/local/bin

RUN git config --global --add safe.directory '*'

COPY entrypoint.sh /entrypoint.sh

WORKDIR /hydrapp/work

CMD /entrypoint.sh
