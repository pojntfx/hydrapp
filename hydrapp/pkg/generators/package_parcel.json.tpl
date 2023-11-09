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
    "@types/react": "^18.0.17",
    "@types/react-dom": "^18.0.6",
    "parcel": "^2.7.0",
    "typescript": "^4.6.4"
  },
  "dependencies": {
    "react": "^18.2.0",
    "react-dom": "^18.2.0",
    "@pojntfx/dudirekta": "^0.6.1"
  }
}
