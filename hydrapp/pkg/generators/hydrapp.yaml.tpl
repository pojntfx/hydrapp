name: hydrapp CI

on:
  push:
  pull_request:
  schedule:
    - cron: "0 0 * * 0"

jobs:
  build-linux:
    runs-on: {{"${{"}} matrix.target.runner {{"}}"}}
    permissions:
      contents: read
    strategy:
      matrix:
        target:
          # {{ .AppID }}
          - id: hydrapp-apk.{{ .AppID }}
            src: .
            pkg: .
            exclude: deb|dmg|flatpak|msi|rpm|binaries|tests
            tag: main
            dst: out/*
            runner: ubuntu-latest
          - id: hydrapp-deb-amd64.{{ .AppID }}
            src: .
            pkg: .
            exclude: deb/arm64|apk|dmg|flatpak|msi|rpm|binaries|tests
            tag: main
            dst: out/*
            runner: ubuntu-latest
          - id: hydrapp-deb-arm64.{{ .AppID }}
            src: .
            pkg: .
            exclude: deb/amd64|apk|dmg|flatpak|msi|rpm|binaries|tests
            tag: main
            dst: out/*
            runner: ubuntu-latest
          - id: hydrapp-dmg.{{ .AppID }}
            src: .
            pkg: .
            exclude: apk|deb|flatpak|msi|rpm|binaries|tests
            tag: main
            dst: out/*
            runner: ubuntu-latest
          - id: hydrapp-flatpak-amd64.{{ .AppID }}
            src: .
            pkg: .
            exclude: apk|deb|flatpak/arm64|dmg|msi|rpm|binaries|tests
            tag: main
            dst: out/*
            runner: ubuntu-latest
          - id: hydrapp-flatpak-arm64.{{ .AppID }}
            src: .
            pkg: .
            exclude: apk|deb|flatpak/amd64|dmg|msi|rpm|binaries|tests
            tag: main
            dst: out/*
            runner: ubuntu-latest
          - id: hydrapp-msi.{{ .AppID }}
            src: .
            pkg: .
            exclude: apk|deb|flatpak|dmg|rpm|binaries|tests
            tag: main
            dst: out/*
            runner: ubuntu-latest
          - id: hydrapp-rpm-amd64.{{ .AppID }}
            src: .
            pkg: .
            exclude: apk|deb|flatpak|dmg|rpm/arm64|msi|binaries|tests
            tag: main
            dst: out/*
            runner: ubuntu-latest
          - id: hydrapp-rpm-arm64.{{ .AppID }}
            src: .
            pkg: .
            exclude: apk|deb|flatpak|dmg|rpm/amd64|msi|binaries|tests
            tag: main
            dst: out/*
            runner: ubuntu-latest
          - id: hydrapp-binaries.{{ .AppID }}
            src: .
            pkg: .
            exclude: apk|deb|flatpak|dmg|msi|rpm|tests
            tag: main
            dst: out/*
            runner: ubuntu-latest
          - id: hydrapp-tests.{{ .AppID }}
            src: .
            pkg: .
            exclude: apk|deb|flatpak|dmg|msi|rpm|binaries
            tag: main
            dst: out/*
            runner: ubuntu-latest

    steps:
      - name: Maximize build space
        run: |
          sudo rm -rf /usr/share/dotnet
          sudo rm -rf /usr/local/lib/android
          sudo rm -rf /opt/ghc
      - name: Checkout
        uses: actions/checkout@v4
      - name: Set up QEMU
        uses: docker/setup-qemu-action@v3
      - name: Set up Docker Buildx
        uses: docker/setup-buildx-action@v3
      - name: Set up hydrapp
        run: |
          curl -L -o /tmp/hydrapp "https://github.com/pojntfx/hydrapp/releases/download/release-main/hydrapp.linux-$(uname -m)" 
          sudo install /tmp/hydrapp /usr/local/bin
      - name: Setup Java/APK keystore
        working-directory: .
        env:
          JAVA_KEYSTORE: {{"${{"}} secrets.JAVA_KEYSTORE {{"}}"}}
        run: echo "${JAVA_KEYSTORE}" | base64 -d >'/tmp/keystore.jks'
      - name: Setup PGP key
        working-directory: .
        env:
          PGP_KEY: {{"${{"}} secrets.PGP_KEY {{"}}"}}
        run: echo "${PGP_KEY}" | base64 -d >'/tmp/pgp.asc'
      - name: Build with hydrapp
        working-directory: {{"${{"}} matrix.target.src {{"}}"}}
        env:
          HYDRAPP_JAVA_KEYSTORE: /tmp/keystore.jks
          HYDRAPP_JAVA_KEYSTORE_PASSWORD: {{"${{"}} secrets.JAVA_KEYSTORE_PASSWORD {{"}}"}}
          HYDRAPP_JAVA_CERTIFICATE_PASSWORD: {{"${{"}} secrets.JAVA_CERTIFICATE_PASSWORD {{"}}"}}

          HYDRAPP_PGP_KEY: /tmp/pgp.asc
          HYDRAPP_PGP_KEY_PASSWORD: {{"${{"}} secrets.PGP_KEY_PASSWORD {{"}}"}}
          HYDRAPP_PGP_KEY_ID: {{"${{"}} secrets.PGP_KEY_ID {{"}}"}}
        run: |
          hydrapp build --config='./{{"${{"}} matrix.target.pkg {{"}}"}}/hydrapp.yaml' --exclude='{{"${{"}} matrix.target.exclude {{"}}"}}' \
            --pull=true --tag='{{"${{"}} matrix.target.tag {{"}}"}}' \
            --dst="${PWD}/out/{{"${{"}} matrix.target.pkg {{"}}"}}" --src="${PWD}" \
            --concurrency="$(nproc)"
      - name: Upload output
        uses: actions/upload-artifact@v4
        with:
          name: {{"${{"}} matrix.target.id {{"}}"}}
          path: {{"${{"}} matrix.target.dst {{"}}"}}

  publish-linux:
    runs-on: ubuntu-latest
    permissions:
      contents: write
      pages: write
      id-token: write
    needs: build-linux
    environment:
      name: github-pages
      {{- if .ExperimentalGithubPagesAction }}
      url: {{"${{"}} steps.publish.outputs.page_url {{"}}"}}
      {{- else }}
      url: {{"${{"}} steps.setup.outputs.base_url {{"}}"}}
      {{- end }}

    steps:
      - name: Checkout
        uses: actions/checkout@v4
      - name: Download output
        uses: actions/download-artifact@v4
        with:
          path: /tmp/out
      - name: Isolate the repositories
        run: |
          mkdir -p /tmp/github-pages
          for dir in /tmp/out/*/; do
            rsync -a "${dir}"/ /tmp/github-pages/
          done

          touch /tmp/github-pages/.nojekyll
      - name: Add index for repositories
        run: |
          sudo apt update
          sudo apt install -y tree

          cd /tmp/github-pages/
          tree --timefmt '%Y-%m-%dT%H:%M:%SZ' -T 'hydrapp Repositories' --du -h -D -H . -o 'index.html'
      {{- if .ExperimentalGithubPagesAction }}
      - name: Setup GitHub Pages
        uses: actions/configure-pages@v5
      - name: Upload GitHub Pages artifact
        uses: actions/upload-pages-artifact@v3
        with:
          path: /tmp/github-pages/
      - name: Publish to GitHub pages
        if: startsWith(github.ref, 'refs/tags/v')
        id: publish
        uses: actions/deploy-pages@v4
      {{- else }}
      - name: Setup GitHub Pages
        id: setup
        uses: actions/configure-pages@v5
      - name: Publish to GitHub pages
        uses: peaceiris/actions-gh-pages@v4
        with:
          github_token: {{"${{"}} secrets.GITHUB_TOKEN {{"}}"}}
          publish_dir: /tmp/github-pages/
          keep_files: true
          user_name: github-actions[bot]
          user_email: github-actions[bot]@users.noreply.github.com
      {{- end }} 
