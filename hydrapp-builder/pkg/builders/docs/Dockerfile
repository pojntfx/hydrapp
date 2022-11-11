FROM fedora:36

# Install native dependencies
RUN dnf install -y pandoc

COPY entrypoint.sh /entrypoint.sh

WORKDIR /work

CMD /entrypoint.sh
