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
make PREFIX=/usr %{?_smp_mflags}

%install
make PREFIX=/usr DESTDIR=%{?buildroot} install
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