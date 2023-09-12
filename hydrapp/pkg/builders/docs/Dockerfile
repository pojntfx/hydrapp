FROM fedora:38

# Install native dependencies
RUN dnf install -y pandoc

COPY entrypoint.sh /entrypoint.sh

WORKDIR /work

CMD /entrypoint.sh
