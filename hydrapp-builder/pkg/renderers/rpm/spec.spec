Name:           {{ .AppID }}
Version:        {{ (LastRelease .AppReleases).Version }}
Release:        1%{?dist}
Summary:        {{ .AppSummary }}

License:        {{ .AppSPDX }}
URL:            {{ .AppURL }}
Source0:        %{name}-%{version}.tar.gz

%if 0%{?suse_version}
BuildRequires: go >= 1.15 ImageMagick >= 6, desktop-file-utils >= 0.23, git >= 2.27.0, appstream-glib >= 0.7.16, npm >= 8.11.0{{ range $pkg := .ExtraPackages }}, {{ $pkg.Name }} >= {{ $pkg.Version }}{{ end }}
%else
BuildRequires: golang >= 1.15 ImageMagick >= 6, desktop-file-utils >= 0.23, git >= 2.27.0, libappstream-glib >= 0.7.14, npm >= 8.11.0{{ range $pkg := .ExtraPackages }}, {{ $pkg.Name }} >= {{ $pkg.Version }}{{ end }}
%endif

Suggests: chromium >= 90

%description
{{ .AppDescription }}

%global debug_package %{nil}

%prep
%autosetup


%build
GOFLAGS="{{ .GoFlags }}" go generate {{ .GoMain }}/...
GOFLAGS="{{ .GoFlags }}" CGO_ENABLED=1 go build -o out/{{ .AppID }} {{ .GoMain }}
for icon in 16x16 22x22 24x24 32x32 36x36 48x48 64x64 72x72 96x96 128x128 192x192 256x256 512x512; do convert "{{ .GoMain }}/icon.png" -resize ${icon} out/icon-${icon}.png; done

%install
install -D out/{{ .AppID }} %{?buildroot}/%{_bindir}/{{ .AppID }}
desktop-file-install --dir=%{?buildroot}/usr/share/applications {{ .GoMain }}/{{ .AppID }}.desktop
appstream-util validate-relax {{ .GoMain }}/{{ .AppID }}.metainfo.xml
for icon in 16x16 22x22 24x24 32x32 36x36 48x48 64x64 72x72 96x96 128x128 192x192 256x256 512x512; do install -D -m 0644 out/icon-${icon}.png %{?buildroot}/usr/share/icons/hicolor/${icon}/apps/{{ .AppID }}.png; done
install -D -m 0644 {{ .GoMain }}/{{ .AppID }}.metainfo.xml ${RPM_BUILD_ROOT}%{_datadir}/metainfo/{{ .AppID }}.metainfo.xml

%files
%license LICENSE
%{_bindir}/{{ .AppID }}
%{_datadir}/applications/{{ .AppID }}.desktop
%{_datadir}/metainfo/{{ .AppID }}.metainfo.xml
%{_datadir}/icons/hicolor/*/apps/{{ .AppID }}.png

%changelog{{ range $release := .AppReleases }}
* {{ $release.Date.Format "Mon Jan 02 2006" }} {{ $release.Author }} <{{ $release.Email }}> {{ $release.Version }}-1
- {{ $release.Description }}
{{ end }}