FROM golang:bullseye

# Install native dependencies
RUN apt update
RUN apt install -y npm

RUN go env -w GOFLAGS=-buildvcs=false
