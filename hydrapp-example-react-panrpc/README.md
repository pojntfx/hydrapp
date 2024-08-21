# hydrapp React and panrpc Example

hydrapp React and panrpc example app.

[![hydrapp CI](https://github.com/pojntfx/hydrapp/actions/workflows/hydrapp.yaml/badge.svg)](https://github.com/pojntfx/hydrapp/actions/workflows/hydrapp.yaml)

## Overview

A simple hydrapp React and panrpc example app.

## Installation

See [INSTALLATION.html](https://pojntfx.github.io/hydrapp/hydrapp-example-react-panrpc//docs/main/INSTALLATION.html).

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

To build and start a development version of hydrapp React and panrpc Example locally, run the following:

```shell
$ git clone https://github.com/pojntfx/hydrapp.git --single-branch
$ cd hydrapp
$ go generate ./hydrapp-example-react-panrpc/...
$ go run ./hydrapp-example-react-panrpc
```

To start the backend and open the frontend in a browser instead of an application window during development, run the following:

```shell
# Start the backend in the first terminal
$ HYDRAPP_BACKEND_LADDR=localhost:1337 HYDRAPP_TYPE=dummy go run ./hydrapp-example-react-panrpc
# Start the frontend in a second terminal
$ cd hydrapp-example-react-panrpc/pkg/frontend
$ npm run dev
# Now open http://localhost:1234 in your browser
```

To build the DEB, RPM, Flatpak, MSI, EXE, DMG, APK, and static binaries for all other platforms, run the following:

```shell
$ go run ./hydrapp build --config ./hydrapp-example-react-panrpc/hydrapp.yaml
# You can find the built packages in the out/ directory
```

If you only want to build certain packages or for certain architectures, for example to only build the APKs, pass `--exclude` like in the following:

```shell
$ go run ./hydrapp build --exclude '(binaries|deb|rpm|flatpak|msi|dmg|docs|tests)' --config ./hydrapp-example-react-panrpc/hydrapp.yaml
```

For more information, see the [hydrapp documentation](../README.md).

## License

hydrapp React and panrpc Example (c) 2024 Felicitas Pojtinger and contributors

SPDX-License-Identifier: Apache-2.0
