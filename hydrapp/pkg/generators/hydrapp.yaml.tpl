name: hydrapp CI

on:
  push:
  pull_request:
  schedule:
    - cron: "0 0 * * 0"

jobs:
  build-linux:
    runs-on: ubuntu-latest
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
          - id: hydrapp-deb.{{ .AppID }}
            src: .
            pkg: .
            exclude: apk|dmg|flatpak|msi|rpm|binaries|tests
            tag: main
            dst: out/*
          - id: hydrapp-dmg.{{ .AppID }}
            src: .
            pkg: .
            exclude: apk|deb|flatpak|msi|rpm|binaries|tests
            tag: main
            dst: out/*
          - id: hydrapp-flatpak.{{ .AppID }}
            src: .
            pkg: .
            exclude: apk|deb|dmg|msi|rpm|binaries|tests
            tag: main
            dst: out/*
          - id: hydrapp-msi.{{ .AppID }}
            src: .
            pkg: .
            exclude: apk|deb|flatpak|dmg|rpm|binaries|tests
            tag: main
            dst: out/*
          - id: hydrapp-rpm.{{ .AppID }}
            src: .
            pkg: .
            exclude: apk|deb|flatpak|dmg|msi|binaries|tests
            tag: main
            dst: out/*
          - id: hydrapp-binaries.{{ .AppID }}
            src: .
            pkg: .
            exclude: apk|deb|flatpak|dmg|msi|rpm|tests
            tag: main
            dst: out/*
          - id: hydrapp-tests.{{ .AppID }}
            src: .
            pkg: .
            exclude: apk|deb|flatpak|dmg|msi|rpm|binaries
            tag: main
            dst: out/*

    steps:
      - name: Maximize build space
        run: |
          sudo rm -rf /usr/share/dotnet
          sudo rm -rf /usr/local/lib/android
          sudo rm -rf /opt/ghc
      - name: Checkout
        uses: actions/checkout@v2
      - name: Set up QEMU
        uses: docker/setup-qemu-action@v1
      - name: Set up Docker Buildx
        uses: docker/setup-buildx-action@v1
      - name: Set up hydrapp
        run: |
          curl -L -o /tmp/hydrapp "https://github.com/pojntfx/hydrapp/releases/latest/download/hydrapp.linux-$(uname -m)" 
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
        working-directory: ${{ matrix.target.src }}
        env:
          HYDRAPP_JAVA_KEYSTORE: /tmp/keystore.jks
          HYDRAPP_JAVA_KEYSTORE_PASSWORD: {{"${{"}} secrets.JAVA_KEYSTORE_PASSWORD {{"}}"}}
          HYDRAPP_JAVA_CERTIFICATE_PASSWORD: {{"${{"}} secrets.JAVA_CERTIFICATE_PASSWORD {{"}}"}}

          HYDRAPP_PGP_KEY: /tmp/pgp.asc
          HYDRAPP_PGP_KEY_PASSWORD: {{"${{"}} secrets.PGP_KEY_PASSWORD {{"}}"}}
          HYDRAPP_PGP_KEY_ID: {{"${{"}} secrets.PGP_KEY_ID {{"}}"}}
        run: |
          export BRANCH_ID=""
          export BRANCH_NAME=""
          if [ "$(git tag --points-at HEAD)" = "" ]; then
            export BRANCH_ID="$(git symbolic-ref --short HEAD)"
            export BRANCH_NAME="$(echo ${BRANCH_ID^})"
          fi

          hydrapp build --config='./{{"${{"}} matrix.target.pkg {{"}}"}}/hydrapp.yaml' --exclude='{{"${{"}} matrix.target.exclude {{"}}"}}' \
            --pull=true --tag='{{"${{"}} matrix.target.tag {{"}}"}}' \
            --dst="${PWD}/out/{{"${{"}} matrix.target.pkg {{"}}"}}" --src="${PWD}" \
            --concurrency="$(nproc)" \
            --branch-id="${BRANCH_ID}" --branch-name="${BRANCH_NAME}"
      - name: Fix permissions for output
        run: sudo chown -R $USER .
      - name: Upload output
        uses: actions/upload-artifact@v2
        with:
          name: {{"${{"}} matrix.target.id {{"}}"}}
          path: {{"${{"}} matrix.target.dst {{"}}"}}

  publish-linux:
    runs-on: ubuntu-latest
    needs: build-linux

    steps:
      - name: Checkout
        uses: actions/checkout@v2
      - name: Download output
        uses: actions/download-artifact@v2
        with:
          path: /tmp/out
      - name: Isolate the repositories
        run: |
          mkdir -p /tmp/github-pages
          for dir in /tmp/out/*/; do
            rsync -a "${dir}"/ /tmp/github-pages/
          done
      - name: Add index for repositories
        run: |
          sudo apt update
          sudo apt install -y tree

          cd /tmp/github-pages/
          tree --timefmt '%Y-%m-%dT%H:%M:%SZ' -T 'hydrapp Repositories' --du -h -D -H . -o 'index.html'
      - name: Publish to GitHub pages
        uses: peaceiris/actions-gh-pages@v3
        with:
          github_token: {{"${{"}} secrets.GITHUB_TOKEN {{"}}"}}
          publish_dir: /tmp/github-pages/
          keep_files: true
          user_name: github-actions[bot]
          user_email: github-actions[bot]@users.noreply.github.com
