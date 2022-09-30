FROM fedora:36

# Install native dependencies
RUN dnf install -y wixl make imagemagick msitools curl wine uuid

# Install MSYS2
RUN curl -L -o /tmp/msys2.exe 'https://github.com/msys2/msys2-installer/releases/download/2021-11-30/msys2-base-x86_64-20211130.sfx.exe'
RUN wine64 /tmp/msys2.exe x -y -oC:/

# Fix MSYS2
RUN sed -i /root/.wine/drive_c/msys64/etc/pacman.conf -e 's/SigLevel    = Required/SigLevel = Never/g'
RUN cat /etc/pki/tls/certs/ca-bundle.crt >/root/.wine/drive_c/msys64/usr/ssl/certs/ca-bundle.crt
RUN cat /etc/pki/tls/certs/ca-bundle.trust.crt >/root/.wine/drive_c/msys64/usr/ssl/certs/ca-bundle.trust.crt

ENV WINEPATH="C:\\msys64\\usr\\bin"

# Install GCC and Go
RUN rm -f /root/.wine/drive_c/msys64/var/lib/pacman/db.lck
RUN wine64 bash.exe -c 'pacman --noconfirm --ignore pacman -Syy'
RUN rm -f /root/.wine/drive_c/msys64/var/lib/pacman/db.lck
RUN wine64 bash.exe -c 'pacman --noconfirm --ignore pacman --needed -S git mingw-w64-x86_64-gcc mingw-w64-x86_64-go mingw-w64-x86_64-pkg-config'
RUN rm -f ~/.wine/drive_c/msys64/var/lib/pacman/db.lck
RUN sed -i /root/.wine/drive_c/msys64/mingw64/lib/pkgconfig/* -e 's/-Wl,-luuid/-luuid/g' # See https://github.com/gotk3/gotk3/wiki/Installing-on-Windows#chocolatey, fails with invalid flag in pkg-config --libs: -Wl,-luuid otherwise
RUN mkdir -p /root/.wine/drive_c/go /root/.wine/drive_c/tmp

COPY entrypoint.sh /entrypoint.sh

WORKDIR /work

CMD /entrypoint.sh
