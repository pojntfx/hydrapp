rm -f *.apk
~/Android/Sdk/build-tools/31.0.0/aapt2 compile -o resources.zip
javac -source 1.8 -target 1.8 -d obj -classpath src -bootclasspath ~/Android/Sdk/platforms/android-31/android.jar src/com/example/helloandroid/*.java
~/Android/Sdk/build-tools/31.0.0/d8 --output . obj/com/example/helloandroid/*
~/Android/Sdk/build-tools/31.0.0/aapt2 link -o android-without-gradle.apk -I ~/Android/Sdk/platforms/android-31/android.jar resources.zip --manifest AndroidManifest.xml
~/Android/Sdk/build-tools/31.0.0/aapt add android-without-gradle.apk classes.dex
[ ! -f mykey.keystore ] && keytool -genkeypair -validity 365 -keystore mykey.keystore -keyalg RSA -keysize 2048 -keypass 123456 -storepass 123456 -dname "cn=Unknown, ou=Unknown, o=Unknown, c=Unknown"
~/Android/Sdk/build-tools/31.0.0/zipalign -f -p 4 android-without-gradle.apk android-without-gradle-output.apk
~/Android/Sdk/build-tools/31.0.0/apksigner sign --ks mykey.keystore --ks-pass pass:123456 --key-pass pass:123456 android-without-gradle-output.apk
adb install android-without-gradle-output.apk
