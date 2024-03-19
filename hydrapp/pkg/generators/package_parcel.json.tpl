{
  "name": "{{ .AppID }}",
  "version": "0.0.1",
  "description": "{{ .AppDescription }}",
  "type": "module",
  "scripts": {
    "dev": "parcel index.html --dist-dir out",
    "build": "tsc && parcel build index.html --dist-dir out"
  },
  "keywords": [],
  "author": "{{ .ReleaseAuthor }} <{{ .ReleaseEmail }}>",
  "license": "{{ .LicenseSPDX }}",
  "devDependencies": {
    "@types/react": "^18.2.67",
    "@types/react-dom": "^18.2.22",
    "parcel": "^2.12.0",
    "process": "^0.11.10",
    "typescript": "^5.4.2"
  },
  "dependencies": {
    "@pojntfx/panrpc": "^0.7.1",
    "@streamparser/json-whatwg": "^0.0.20",
    "react": "^18.2.0",
    "react-dom": "^18.2.0",
    "use-async": "^1.1.0"
  },
  "@parcel/resolver-default": {
    "packageExports": true
  }
}
