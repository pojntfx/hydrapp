Name:           desktop-integrated-webserver-rpm
Version:        0.0.1
Release:        1%{?dist}
Summary:        Simple hello world example

License:        AGPLv3
URL:            https://github.com/pojntfx/multi-browser-electron
Source0:        %{name}-%{version}.tar.gz

BuildRequires:  golang >= 1.15

Suggests: chromium >= 90

%description
A simple hello world example to demonstrate RPM packaging.

%prep
%autosetup


%build
make PREFIX=/usr %{?_smp_mflags}

%install
make PREFIX=/usr DESTDIR=%{?buildroot} install
mkdir -p $RPM_BUILD_ROOT/%{_mandir}/man1/
mv docs/desktop-integrated-webserver-rpm.1 $RPM_BUILD_ROOT/%{_mandir}/man1/desktop-integrated-webserver-rpm.1

%files
%license LICENSE
%{_bindir}/desktop-integrated-webserver-rpm
%{_mandir}/man1/desktop-integrated-webserver-rpm.1*


%changelog
* Mon Aug 30 2021 Felicitas Pojtinger <felicitas@pojtinger.com> 0.0.1-1
- Initial release.