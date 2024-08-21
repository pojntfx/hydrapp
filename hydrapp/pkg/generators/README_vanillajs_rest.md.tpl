# {{ .AppName }}

{{ .AppSummary }}.

[![hydrapp CI]({{ .AppGitWeb }}/actions/workflows/hydrapp.yaml/badge.svg)]({{ .AppGitWeb }}/actions/workflows/hydrapp.yaml)

## Overview

{{ .AppDescription }}

## Installation

See [INSTALLATION.html]({{ .AppBaseURL }}/docs/main/INSTALLATION.html).

## Reference

### Command Line Arguments

All arguments passed to the binary will be forwarded to the browser used to display the frontend.

### Environment Variables

| Name                     | Description                                                                                                 |
| ------------------------ | ----------------------------------------------------------------------------------------------------------- |
| `HYDRAPP_BACKEND_LADDR`  | Listen address for the backend (`localhost:0` by default)                                                   |
| `HYDRAPP_FRONTEND_LADDR` | Listen address for the frontend (`localhost:0` by default)                                                  |
| `HYDRAPP_BROWSER`        | Binary of browser to display the frontend with                                                              |
| `HYDRAPP_TYPE`           | Type of browser to display the frontend with (one of `chromium`, `firefox`, `epiphany`, `lynx` and `dummy`) |
| `HYDRAPP_SELFUPDATE`     | Whether to check for updates on launch (disabled if OS provides an app update mechanism)                    |

## Acknowledgements

- [pojntfx/hydrapp](https://github.com/pojntfx/hydrapp) provides the application framework.

## Contributing

To contribute, please use the [GitHub flow](https://guides.github.com/introduction/flow/) and follow our [Code of Conduct](./CODE_OF_CONDUCT.md).

To build and start a development version of {{ .AppName }} locally, run the following:

```shell
{{- if .ExperimentalGithubPagesAction }}
$ git clone {{ .AppGit }} --single-branch
{{- else }}
$ git clone {{ .AppGit }}
{{- end }}
$ cd {{ .Dir }}
$ go generate ./...
$ go run .
```

To build the DEB, RPM, Flatpak, MSI, EXE, DMG, APK, and static binaries for all other platforms, run the following:

```shell
$ hydrapp build
# You can find the built packages in the out/ directory
```

If you only want to build certain packages or for certain architectures, for example to only build the APKs, pass `--exclude` like in the following:

```shell
$ hydrapp build --exclude '(binaries|deb|rpm|flatpak|msi|dmg|docs|tests)'
```

For more information, see the [hydrapp documentation](https://github.com/pojntfx/hydrapp).

## License

{{ .AppName }} (c) {{ .CurrentYear }} {{ .ReleaseAuthor }} and contributors

SPDX-License-Identifier: {{ .LicenseSPDX }}
