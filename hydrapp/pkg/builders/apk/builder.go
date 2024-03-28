package apk

import (
	"context"
	"encoding/base64"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/docker/docker/client"
	"github.com/pojntfx/hydrapp/hydrapp/pkg/builders"
	"github.com/pojntfx/hydrapp/hydrapp/pkg/executors"
	"github.com/pojntfx/hydrapp/hydrapp/pkg/renderers"
	"github.com/pojntfx/hydrapp/hydrapp/pkg/renderers/apk"
	"github.com/pojntfx/hydrapp/hydrapp/pkg/utils"
)

const (
	Image = "ghcr.io/pojntfx/hydrapp-build-apk"
)

func NewBuilder(
	ctx context.Context,
	cli *client.Client,

	image string, // OCI image to use
	pull bool, // Whether to pull the image or not
	src, // Input directory
	dst string, // Output directory
	onID func(id string), // Callback to handle container ID
	onOutput func(shortID string, color string, timestamp int64, message string), // Callback to handle container output
	appID string, // Android app ID to use
	javaKeystore []byte, // Android cert contents
	javaKeystorePassword string, // Password for the Android keystore
	javaCertificatePassword string, // Password for the Android certificate
	pgpKey []byte, // PGP key contents
	pgpKeyPassword string, // Password for the PGP key
	baseURL, // Base URL where the repo is to be hosted
	appName string, // App name
	overwrite bool, // Overwrite files even if they exist
	branchID, // Branch ID
	branchName, // Branch Name
	goMain, // Directory with the main package to build
	goFlags, // Flags to pass to the Go command
	goGenerate string, // Command to execute go generate with
) *Builder {
	return &Builder{
		ctx,
		cli,

		image,
		pull,
		src,
		dst,
		onID,
		onOutput,
		appID,
		base64.StdEncoding.EncodeToString(javaKeystore),
		javaKeystorePassword,
		javaCertificatePassword,
		base64.StdEncoding.EncodeToString(pgpKey),
		pgpKeyPassword,
		baseURL,
		appName,
		overwrite,
		branchID,
		branchName,
		goMain,
		goFlags,
		goGenerate,
	}
}

type Builder struct {
	ctx context.Context
	cli *client.Client

	image string
	pull  bool
	src,
	dst string
	onID     func(id string)
	onOutput func(shortID string, color string, timestamp int64, message string)
	appID,
	javaKeystore,
	javaKeystorePassword,
	javaCertificatePassword,
	pgpKey,
	pgpKeyPassword,
	baseURL,
	appName string
	overwrite bool
	branchID,
	branchName,
	goMain,
	goFlags,
	goGenerate string
}

func (b *Builder) Render(workdir string, ejecting bool) error {
	appID := builders.GetAppIDForBranch(b.appID, b.branchID)
	appName := builders.GetAppNameForBranch(b.appName, b.branchName)

	if strings.TrimSpace(b.branchID) != "" {
		jniBindingsPath := filepath.Join(workdir, b.goMain, "android.go")

		stableJNIBindingsContent, err := os.ReadFile(jniBindingsPath)
		if err != nil {
			return err
		}

		stableJavaID := strings.Replace(b.appID, ".", "_", -1)
		mainJavaID := strings.Replace(appID, ".", "_", -1)

		if !ejecting || b.overwrite {
			if !strings.Contains(string(stableJNIBindingsContent), mainJavaID) {
				if err := ioutil.WriteFile(jniBindingsPath, []byte(strings.Replace(string(stableJNIBindingsContent), stableJavaID, mainJavaID, -1)), 0664); err != nil {
					return err
				}
			}
		}
	}

	return utils.WriteRenders(
		filepath.Join(workdir, b.goMain),
		[]*renderers.Renderer{
			apk.NewManifestRenderer(
				appID,
				appName,
			),
			apk.NewActivityRenderer(
				appID,
			),
			apk.NewHeaderRenderer(),
			apk.NewImplementationRenderer(),
		},
		b.overwrite,
		ejecting,
	)
}

func (b *Builder) Build() error {
	dst := builders.GetFilepathForBranch(b.dst, b.branchID)
	appID := builders.GetAppIDForBranch(b.appID, b.branchID)
	baseURL := builders.GetPathForBranch(b.baseURL, b.branchID, "") + "/repo" // F-Droid requires the path to end with `/repo`: `CRITICAL: repo_url needs to end with /repo`

	return executors.DockerRunImage(
		b.ctx,
		b.cli,

		b.image,
		b.pull,
		false,
		b.src,
		dst,
		b.onID,
		b.onOutput,
		map[string]string{
			"APP_ID":                    appID,
			"JAVA_KEYSTORE":             b.javaKeystore,
			"JAVA_KEYSTORE_PASSWORD":    b.javaKeystorePassword,
			"JAVA_CERTIFICATE_PASSWORD": b.javaCertificatePassword,
			"PGP_KEY":                   b.pgpKey,
			"PGP_KEY_PASSWORD":          b.pgpKeyPassword,
			"BASE_URL":                  baseURL,
			"GOMAIN":                    b.goMain,
			"GOFLAGS":                   b.goFlags,
			"GOGENERATE":                b.goGenerate,
		},
		b.Render,
		[]string{},
	)
}
