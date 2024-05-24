FROM alpine

# Install native dependencies
RUN apk add flatpak-builder xz bash appstream-compose 

COPY entrypoint.sh /entrypoint.sh

WORKDIR /hydrapp/work

CMD /entrypoint.sh
