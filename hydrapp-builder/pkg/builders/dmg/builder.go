package dmg

import (
	"context"
	"encoding/base64"
	"path/filepath"
	"strings"

	"github.com/docker/docker/client"
	"github.com/pojntfx/hydrapp/hydrapp-builder/pkg/builders"
	"github.com/pojntfx/hydrapp/hydrapp-builder/pkg/executors"
	"github.com/pojntfx/hydrapp/hydrapp-builder/pkg/renderers"
	"github.com/pojntfx/hydrapp/hydrapp-builder/pkg/renderers/dmg"
	"github.com/pojntfx/hydrapp/hydrapp-builder/pkg/utils"
)

const (
	Image = "ghcr.io/pojntfx/hydrapp-build-dmg"
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
	appID, // macOS app ID to use
	appName string, // Human-readable name for the app
	gpgKeyContent []byte, // GPG key contents
	gpgKeyPassword string, // Password for the GPG key
	universal bool, // Build universal (amd64 and arm64) binary instead of amd64 only
	packages []string, // MacPorts packages to install
	releases []renderers.Release, // App releases
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
		appName,
		base64.StdEncoding.EncodeToString(gpgKeyContent),
		base64.StdEncoding.EncodeToString([]byte(gpgKeyPassword)),
		universal,
		packages,
		releases,
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
	appName,
	gpgKeyContent,
	gpgKeyPassword string
	universal bool
	packages  []string
	releases  []renderers.Release
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

	return utils.WriteRenders(
		filepath.Join(workdir, b.goMain),
		[]*renderers.Renderer{
			dmg.NewInfoRenderer(
				appID,
				appName,
				b.releases,
			),
		},
		b.overwrite,
		ejecting,
	)
}

func (b *Builder) Build() error {
	dst := builders.GetFilepathForBranch(b.dst, b.branchID)
	appID := builders.GetAppIDForBranch(b.appID, b.branchID)
	appName := builders.GetAppNameForBranch(b.appName, b.branchName)

	return executors.DockerRunImage(
		b.ctx,
		b.cli,
		b.image,
		b.pull,
		true,
		b.src,
		dst,
		b.onID,
		b.onOutput,
		map[string]string{
			"APP_ID":           appID,
			"APP_NAME":         appName,
			"GPG_KEY_CONTENT":  b.gpgKeyContent,
			"GPG_KEY_PASSWORD": b.gpgKeyPassword,
			"ARCHITECTURES": func() string {
				if b.universal {
					return "amd64 arm64"
				}

				return "amd64"
			}(),
			"MACPORTS":   strings.Join(b.packages, " "),
			"GOMAIN":     b.goMain,
			"GOFLAGS":    b.goFlags,
			"GOGENERATE": b.goGenerate,
		},
		b.Render,
	)
}
