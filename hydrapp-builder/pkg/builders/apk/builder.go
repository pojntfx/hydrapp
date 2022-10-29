package apk

import (
	"context"
	"encoding/base64"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/docker/docker/client"
	"github.com/pojntfx/hydrapp/hydrapp-builder/pkg/builders"
	"github.com/pojntfx/hydrapp/hydrapp-builder/pkg/executors"
	"github.com/pojntfx/hydrapp/hydrapp-builder/pkg/renderers"
	"github.com/pojntfx/hydrapp/hydrapp-builder/pkg/renderers/apk"
	"github.com/pojntfx/hydrapp/hydrapp-builder/pkg/utils"
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
	pgpKeyContent []byte, // PGP key contents
	pgpKeyPassword string, // Password for the PGP key
	androidCertContent []byte, // Android cert contents
	androidStorepass string, // Password for the Android keystore
	androidKeypass string, // Password for the Android certificate
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
		base64.StdEncoding.EncodeToString(pgpKeyContent),
		base64.StdEncoding.EncodeToString([]byte(pgpKeyPassword)),
		base64.StdEncoding.EncodeToString(androidCertContent),
		base64.StdEncoding.EncodeToString([]byte(androidStorepass)),
		base64.StdEncoding.EncodeToString([]byte(androidKeypass)),
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
	pgpKeyContent,
	pgpKeyPassword,
	androidCertContent,
	androidStorepass,
	androidKeypass,
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

		stableJNIBindingsContent, err := ioutil.ReadFile(jniBindingsPath)
		if err != nil {
			return err
		}

		stableJavaID := strings.Replace(b.appID, ".", "_", -1)
		unstableJavaID := strings.Replace(appID, ".", "_", -1)

		if !ejecting || b.overwrite {
			if !strings.Contains(string(stableJNIBindingsContent), unstableJavaID) {
				if err := ioutil.WriteFile(jniBindingsPath, []byte(strings.Replace(string(stableJNIBindingsContent), stableJavaID, unstableJavaID, -1)), os.ModePerm); err != nil {
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
	baseURL := builders.GetPathForBranch(b.baseURL, b.branchID)

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
			"APP_ID":               appID,
			"PGP_KEY_CONTENT":      b.pgpKeyContent,
			"PGP_KEY_PASSWORD":     b.pgpKeyPassword,
			"ANDROID_CERT_CONTENT": b.androidCertContent,
			"ANDROID_STOREPASS":    b.androidStorepass,
			"ANDROID_KEYPASS":      b.androidKeypass,
			"BASE_URL":             baseURL,
			"GOMAIN":               b.goMain,
			"GOFLAGS":              b.goFlags,
			"GOGENERATE":           b.goGenerate,
		},
		b.Render,
	)
}
