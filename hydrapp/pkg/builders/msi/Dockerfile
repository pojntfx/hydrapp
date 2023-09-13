FROM scottyhardy/docker-wine

RUN apt update
RUN apt install -y curl

# Install MSYS2
RUN curl -L -o /tmp/msys2.exe 'https://github.com/msys2/msys2-installer/releases/download/2023-07-18/msys2-base-x86_64-20230718.sfx.exe'
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
RUN mkdir -p '/root/.wine/drive_c/go' '/root/.wine/drive_c/tmp'
RUN mkdir -p '/root/.wine/drive_c/users/root/Documents/go-workspace'

RUN yes | cpan install File::Slurp

RUN git config --global --add safe.directory '*'

COPY entrypoint.sh /entrypoint.sh

WORKDIR /work

CMD /entrypoint.sh
