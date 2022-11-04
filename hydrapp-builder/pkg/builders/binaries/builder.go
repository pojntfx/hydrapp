package binaries

import (
	"context"
	"encoding/base64"
	"path/filepath"
	"strings"

	"github.com/docker/docker/client"
	"github.com/pojntfx/hydrapp/hydrapp-builder/pkg/builders"
	"github.com/pojntfx/hydrapp/hydrapp-builder/pkg/executors"
	"github.com/pojntfx/hydrapp/hydrapp-builder/pkg/renderers"
	"github.com/pojntfx/hydrapp/hydrapp-builder/pkg/utils"
)

const (
	Image = "ghcr.io/pojntfx/hydrapp-build-binaries"
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
	appID string, // App ID to use
	pgpKeyContent []byte, // PGP key contents
	pgpKeyPassword, // password for the PGP key
	appName string, // App name
	branchID, // Branch ID
	branchName, // Branch Name
	goMain, // Directory with the main package to build
	goFlags, // Flags to pass to the Go command
	goGenerate, // Command to execute go generate with
	goExclude string, // Regex of platforms to ignore
	hostPackages []string, // Debian packages to install before building
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
		appName,
		branchID,
		branchName,
		goMain,
		goFlags,
		goGenerate,
		goExclude,
		hostPackages,
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
	appName,
	branchID,
	branchName,
	goMain,
	goFlags,
	goGenerate,
	goExclude string
	hostPackages []string
}

func (b *Builder) Render(workdir string, ejecting bool) error {
	return utils.WriteRenders(
		filepath.Join(workdir, b.goMain),
		[]*renderers.Renderer{},
		false,
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
			"PGP_KEY_CONTENT":  b.pgpKeyContent,
			"PGP_KEY_PASSWORD": b.pgpKeyPassword,
			"APP_NAME":         appName,
			"GOMAIN":           b.goMain,
			"GOFLAGS":          b.goFlags,
			"GOGENERATE":       b.goGenerate,
			"GOEXCLUDE":        b.goExclude,
			"HOST_PACKAGES":    strings.Join(b.hostPackages, " "),
		},
		b.Render,
	)
}
