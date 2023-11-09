{
  "name": "{{ .AppID }}",
  "version": "0.0.1",
  "description": "{{ .AppDescription }}",
  "type": "module",
  "scripts": {
    "dev": "BUILD_PATH=out react-scripts start",
    "build": "BUILD_PATH=out react-scripts build; cp ../../icon.png out/icon.png"
  },
  "keywords": [],
  "author": "{{ .ReleaseAuthor }} <{{ .ReleaseEmail }}>",
  "license": "{{ .LicenseSPDX }}",
  "devDependencies": {
    "@types/react": "^18.0.17",
    "@types/react-dom": "^18.0.6",
    "react-scripts": "5.0.1",
    "typescript": "^4.6.4"
  },
  "dependencies": {
    "react": "^18.2.0",
    "react-dom": "^18.2.0",
    "@pojntfx/dudirekta": "^0.6.1"
  },
  "eslintConfig": {
    "extends": [
      "react-app",
      "react-app/jest"
    ]
  },
  "browserslist": {
    "production": [
      ">0.2%",
      "not dead",
      "not op_mini all"
    ],
    "development": [
      "last 1 chrome version",
      "last 1 firefox version",
      "last 1 safari version"
    ]
  }
}
