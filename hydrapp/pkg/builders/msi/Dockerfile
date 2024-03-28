FROM golang:bookworm

RUN apt update
RUN apt install -y curl make imagemagick msitools curl uuid-runtime git cpanminus npm tree jq apt-transport-https ca-certificates cabextract git gnupg gosu gpg-agent locales p7zip pulseaudio pulseaudio-utils sudo tzdata unzip wget winbind xvfb zenity wixl libxml2-utils

# Install WINE build dependencies
RUN sed -i 's/^Types: deb$/Types: deb deb-src/' /etc/apt/sources.list.d/debian.sources
RUN apt update
RUN DEBIAN_FRONTEND="noninteractive" apt-get build-dep --install-recommends -y wine

# Build and install patched WINE (see https://github.com/msys2/MSYS2-packages/issues/682)
RUN git clone https://gitlab.winehq.org/jhol/wine.git /tmp/winesrc
WORKDIR /tmp/winesrc
RUN git checkout msys2-hacks-17
RUN ./configure --enable-win64
RUN make -j $(nproc)
RUN make install
RUN rm -rf /tmp/winesrc

# Install MSYS2
RUN curl -L -o /tmp/msys2.exe 'https://github.com/msys2/msys2-installer/releases/download/2024-01-13/msys2-base-x86_64-20240113.sfx.exe'
RUN wine64 /tmp/msys2.exe x -y -oC:/

# Fix MSYS2
RUN sed -i /root/.wine/drive_c/msys64/etc/pacman.conf -e 's/SigLevel    = Required/SigLevel = Never/g'

ENV WINEPATH="C:\\ucrt64\\bin;C:\\msys64\\usr\\bin"
ENV MSYSTEM="UCRT64"

# Install GCC and Go
RUN yes | wine64 bash.exe -c 'pacman --noconfirm --ignore pacman -Syu --debug'

COPY perl-Convert-BinHex-999-1-x86_64.pkg.tar.gz /root/.wine/drive_c/msys64/tmp

RUN wine64 bash.exe -c 'pacman --debug --noconfirm -U /tmp/perl-Convert-BinHex-999-1-x86_64.pkg.tar.gz'
RUN wine64 bash.exe -c 'pacman --debug --noconfirm --ignore pacman --needed -S ca-certificates git mingw-w64-ucrt-x86_64-gcc mingw-w64-ucrt-x86_64-go mingw-w64-ucrt-x86_64-pkg-config'
RUN wine64 bash.exe -c 'update-ca-trust'
RUN sed -i /root/.wine/drive_c/msys64/ucrt64/lib/pkgconfig/* -e 's/-Wl,-luuid/-luuid/g' # See https://github.com/gotk3/gotk3/wiki/Installing-on-Windows#chocolatey, fails with invalid flag in pkg-config --libs: -Wl,-luuid otherwise

RUN yes | cpan install File::Slurp

RUN git config --global --add safe.directory '*'

RUN mv /root/.wine/drive_c/msys64 /tmp/msys64

COPY entrypoint.sh /entrypoint.sh

WORKDIR /work

CMD /entrypoint.sh
