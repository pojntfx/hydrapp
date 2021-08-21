package main

import (
	"flag"
	"io"
	"log"
	"os"
	"path/filepath"
)

func main() {
	// Define the input flags
	// androidAPIVersion := flag.Int("androidAPIVersion", 31, "Version of the Android API to build for")
	// androidBuildToolsVersion := flag.String("androidBuildToolsVersion", "31.0.0", "Version of the Android build tools to use")
	// androidCertificatePassword := flag.String("androidCertificatePassword", "123456", "The password to use for the signing certificate")
	// androidAppID := flag.String("androidAppID", "com.example.app", "The Android app ID, in reverse domain notation.")
	// dist := flag.String("dist", "out", "Directory to build into")

	// Parse the input flags
	flag.Parse()

	// Create the working directory
	workdir, err := os.MkdirTemp(os.TempDir(), "hydrapp-*")
	if err != nil {
		log.Fatalln("could not create work directory:", err)
	}

	// Copy MainActivity.java to the working directory
	if err := copyFile("MainActivity.java", filepath.Join(workdir, "MainActivity.java")); err != nil {
		log.Fatalln("could not copy MainActivity.java:", err)
	}
}

func copyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.OpenFile(dst, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(srcFile, dstFile)

	return err
}
