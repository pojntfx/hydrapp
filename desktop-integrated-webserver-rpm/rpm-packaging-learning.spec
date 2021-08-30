Name:           desktop-integrated-webserver-rpm
Version:        0.0.1
Release:        1%{?dist}
Summary:        simple hello world example

License:        AGPLv3
URL:            https://github.com/pojntfx/rpm-packaging-learning
Source0:        %{name}-%{version}.tar.gz

BuildRequires:  golang

Provides:       %{name} = %{version}

%description
A simple hello world example to demonstrate RPM packaging.

%prep
%autosetup


%build
make PREFIX=/usr %{?_smp_mflags}

%install
make PREFIX=/usr DESTDIR=%{?buildroot} install


%files
%license LICENSE
%{_bindir}/${name}


%changelog
* Mon Aug 30 2021 Felicitas Pojtinger <felicitas@pojtinger.com>
- First release%changelog
