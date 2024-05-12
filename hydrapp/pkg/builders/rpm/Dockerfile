FROM fedora:40

# Install native dependencies
RUN dnf install -y fedora-packager @development-tools qemu-user-static rpm-sign

RUN git config --global --add safe.directory '*'

COPY entrypoint.sh /entrypoint.sh

WORKDIR /hydrapp/work

CMD /entrypoint.sh
