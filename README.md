# hydrapp

![Logo](./docs/logo-readme.png)

Build apps that run everywhere with Go and a browser engine of your choice.

[![hydrapp CI](https://github.com/pojntfx/hydrapp/actions/workflows/hydrapp.yaml/badge.svg)](https://github.com/pojntfx/hydrapp/actions/workflows/hydrapp.yaml)
![Go Version](https://img.shields.io/badge/go%20version-%3E=1.19-61CFDD.svg)
[![Go Reference](https://pkg.go.dev/badge/github.com/pojntfx/hydrapp/hydrapp.svg)](https://pkg.go.dev/github.com/pojntfx/hydrapp/hydrapp)
[![Matrix](https://img.shields.io/matrix/hydrapp:matrix.org)](https://matrix.to/#/#hydrapp:matrix.org?via=matrix.org)

## Overview

hydrapp is a Go framework similar to Electron with unique feature: It **can use (almost) any browser engine** to render the frontend!

It enables you too ...

- **Build apps in Go and JS:** Use the speedy and easy-to-learn Go language to create your app's backend, then use your web tech know-how to develop a top-notch, user-friendly frontend.
- **Connect frontend and backend with ease:** With hydrapp and [dudirekta](https://github.com/pojntfx/dudirekta), you can easily call functions between the frontend and backend without any complicated manual setup.
- **Compatible with all browsers:** Hydrapp works with any pre-installed browser by starting it in PWA mode, so you can render your app on Chrome, Firefox/Gecko, Epiphany/Webkit, and even Android WebView.
- **Cross-compile with CGo easily:** Hydrapp simplifies cross-compiling with a container-based environment that includes MacPorts, MSYS2, APT, and DNF, making it easy to distribute binaries without using non-Linux machines.
- **Effortlessly build, sign, distribute, and update your app:** Hydrapp streamlines your app's delivery with an integrated CI/CD workflow, producing reproducible packages for DEB, RPM, Flatpak, MSI, EXE, DMG, APK, and static binaries for all other platforms. Hydrapp can also generate APT, YUM, and Flatpak repositories for Linux and F-Droid repositories for Android. Self-updating for Windows, macOS, and other platforms is also available.

🚧 This project is a work-in-progress! Instructions will be added as soon as it is usable. 🚧

## License

hydrapp (c) 2023 Felicitas Pojtinger and contributors

SPDX-License-Identifier: AGPL-3.0
