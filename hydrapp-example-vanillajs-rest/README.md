# hydrapp Vanilla JS and REST

hydrapp Vanilla JS and REST app.

[![hydrapp CI](https://github.com/pojntfx/hydrapp/actions/workflows/hydrapp.yaml/badge.svg)](https://github.com/pojntfx/hydrapp/actions/workflows/hydrapp.yaml)

## Overview

A simple hydrapp Vanilla JS and REST app.

## Installation

See [INSTALLATION.html](https://pojntfx.github.io/hydrapp/hydrapp-example-vanillajs-rest//docs/main/INSTALLATION.html).

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

To build and start a development version of hydrapp Vanilla JS and REST locally, first install [Go](https://go.dev/), then run the following:

```shell
$ git clone https://github.com/pojntfx/hydrapp.git --single-branch
$ cd hydrapp
$ go generate ./hydrapp-example-vanillajs-rest/...
$ go run ./hydrapp-example-vanillajs-rest
```

To build the DEB, RPM, Flatpak, MSI, EXE, DMG, APK, and static binaries for all other platforms, run the following:

```shell
$ go run ./hydrapp build --config ./hydrapp-example-vanillajs-rest/hydrapp.yaml
# You can find the built packages in the out/ directory
```

If you only want to build certain packages or for certain architectures, for example to only build the APKs, pass `--exclude` like in the following:

```shell
$ go run ./hydrapp build --exclude '(binaries|deb|rpm|flatpak|msi|dmg|docs|tests)' --config ./hydrapp-example-vanillajs-rest/hydrapp.yaml
```

For more information, see the [hydrapp documentation](../README.md).

## License

hydrapp Vanilla JS and REST (c) 2025 Felicitas Pojtinger and contributors

SPDX-License-Identifier: Apache-2.0
