Name:           com.pojtinger.felicitas.hydrappexample
Version:        0.0.1
Release:        1%{?dist}
Summary:        hydrapp example app

License:        AGPLv3
URL:            https://github.com/pojntfx/multi-browser-electron
Source0:        %{name}-%{version}.tar.gz

%if 0%{?suse_version}
BuildRequires: go >= 1.15 ImageMagick >= 6, desktop-file-utils >= 0.23, git >= 2.27.0
%else
BuildRequires: golang >= 1.15 ImageMagick >= 6, desktop-file-utils >= 0.23, git >= 2.27.0
%endif

Suggests: chromium >= 90

%description
A simple Hydrapp example app.

%global debug_package %{nil}

%prep
%autosetup


%build
for icon in 16x16 22x22 24x24 32x32 36x36 48x48 64x64 72x72 96x96 128x128 192x192 256x256 512x512; do convert icon.png -resize ${icon} out/icon-${icon}.png; done
CGO_ENABLED=1 make PREFIX=/usr %{?_smp_mflags}

%install
CGO_ENABLED=1 make PREFIX=/usr DESTDIR=%{?buildroot} install
desktop-file-install --dir=%{?buildroot}/usr/share/applications com.pojtinger.felicitas.hydrappexample.desktop
for icon in 16x16 22x22 24x24 32x32 36x36 48x48 64x64 72x72 96x96 128x128 192x192 256x256 512x512; do install -D -m 0644 out/icon-${icon}.png %{?buildroot}/usr/share/icons/hicolor/${icon}/apps/com.pojtinger.felicitas.hydrappexample.png; done
install -D -m 0644 com.pojtinger.felicitas.hydrappexample.metainfo.xml ${RPM_BUILD_ROOT}%{_datadir}/metainfo/com.pojtinger.felicitas.hydrappexample.metainfo.xml
install -D -m 0644 docs/com.pojtinger.felicitas.hydrappexample.1 $RPM_BUILD_ROOT/%{_mandir}/man1/com.pojtinger.felicitas.hydrappexample.1

%files
%license LICENSE
%{_bindir}/com.pojtinger.felicitas.hydrappexample
%{_mandir}/man1/com.pojtinger.felicitas.hydrappexample.1*
%{_datadir}/applications/com.pojtinger.felicitas.hydrappexample.desktop
%{_datadir}/metainfo/com.pojtinger.felicitas.hydrappexample.metainfo.xml
%{_datadir}/icons/hicolor/*/apps/com.pojtinger.felicitas.hydrappexample.png


%changelog
* Mon Aug 30 2021 Felicitas Pojtinger <felicitas@pojtinger.com> 0.0.1-1
- Initial release.