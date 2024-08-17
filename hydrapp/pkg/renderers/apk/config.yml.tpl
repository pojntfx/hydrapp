---
sdk_path: "{ ANDROID_HOME }"
repo_url: "{ BASE_URL }"
repo_name: {{ .AppName }} F-Droid Repo
repo_description: >-
  Android apps for {{ .AppName }}.
repo_icon: icon.png
repo_keyalias: { ANDROID_CERT_ALIAS }
keystore: keystore.p12
keystorepass: "{ JAVA_KEYSTORE_PASSWORD }"
keypass: "{ JAVA_CERTIFICATE_PASSWORD }"
keydname: CN={ ANDROID_CERT_CN }
apksigner: /usr/bin/apksigner