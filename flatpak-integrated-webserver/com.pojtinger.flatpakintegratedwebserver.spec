Name:           com.pojtinger.flatpakintegratedwebserver
Version:        0.0.1
Release:        1%{?dist}
Summary:        Simple hello world example

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
A simple hello world example to demonstrate packaging.

%global debug_package %{nil}

%prep
%autosetup


%build
make PREFIX=/usr %{?_smp_mflags}

%install
make PREFIX=/usr DESTDIR=%{?buildroot} install
install -D -m 0644 com.pojtinger.flatpakintegratedwebserver.metainfo.xml ${RPM_BUILD_ROOT}%{_datadir}/metainfo/com.pojtinger.flatpakintegratedwebserver.metainfo.xml
install -D -m 0644 docs/com.pojtinger.flatpakintegratedwebserver.1 $RPM_BUILD_ROOT/%{_mandir}/man1/com.pojtinger.flatpakintegratedwebserver.1

%files
%license LICENSE
%{_bindir}/com.pojtinger.flatpakintegratedwebserver
%{_mandir}/man1/com.pojtinger.flatpakintegratedwebserver.1*
%{_datadir}/applications/com.pojtinger.flatpakintegratedwebserver.desktop
%{_datadir}/metainfo/com.pojtinger.flatpakintegratedwebserver.metainfo.xml
%{_datadir}/icons/hicolor/*/apps/com.pojtinger.flatpakintegratedwebserver.png


%changelog
* Mon Aug 30 2021 Felicitas Pojtinger <felicitas@pojtinger.com> 0.0.1-1
- Initial release.