FROM fedora:39

# Install native dependencies
RUN dnf install -y fedora-packager @development-tools qemu-user-static rpm-sign

COPY entrypoint.sh /entrypoint.sh

WORKDIR /work

CMD /entrypoint.sh
