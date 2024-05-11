package msi

import (
	"context"
	"encoding/base64"
	"io"
	"path/filepath"
	"strings"

	"github.com/docker/docker/client"
	"github.com/pojntfx/hydrapp/hydrapp/pkg/builders"
	"github.com/pojntfx/hydrapp/hydrapp/pkg/executors"
	"github.com/pojntfx/hydrapp/hydrapp/pkg/renderers"
	"github.com/pojntfx/hydrapp/hydrapp/pkg/renderers/msi"
	"github.com/pojntfx/hydrapp/hydrapp/pkg/utils"
)

const (
	Image = "ghcr.io/pojntfx/hydrapp-build-msi"
)

func NewBuilder(
	ctx context.Context,
	cli *client.Client,

	image string, // OCI image to use
	pull bool, // Whether to pull the image or not
	src, // Input directory
	dst string, // Output directory
	onID func(id string), // Callback to handle container ID
	stdout io.Writer, // Writer to handle container output
	appID, // Android app ID to use
	appName string, // Human-readable name for the app,
	pgpKey []byte, // PGP key contents
	pgpKeyPassword, // Password for the PGP key
	architecture string, // Architecture to build for
	packages []string, // MSYS2 packages to install. Only supported for amd64.
	releases []renderers.Release, // App releases
	overwrite bool, // Overwrite files even if they exist
	branchID, // Branch ID
	branchName, // Branch Name
	goMain, // Directory with the main package to build
	goFlags, // Flags to pass to the Go command
	include, // Regex of files and directories from MSYS2 to include
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
		stdout,
		appID,
		appName,
		base64.StdEncoding.EncodeToString(pgpKey),
		pgpKeyPassword,
		architecture,
		packages,
		releases,
		overwrite,
		branchID,
		branchName,
		goMain,
		goFlags,
		include,
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
	onID   func(id string)
	stdout io.Writer
	appID,
	appName,
	pgpKey,
	pgpKeyPassword,
	architecture string
	packages  []string
	releases  []renderers.Release
	overwrite bool
	branchID,
	branchName,
	goMain,
	goFlags,
	include,
	goGenerate string
}

func (b *Builder) Render(workdir string, ejecting bool) error {
	appID := builders.GetAppIDForBranch(b.appID, b.branchID)
	appName := builders.GetAppNameForBranch(b.appName, b.branchName)

	return utils.WriteRenders(
		filepath.Join(workdir, b.goMain),
		[]renderers.Renderer{
			xdg.NewIconRenderer(
			msi.NewWixRenderer(
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
		false,
		b.src,
		dst,
		b.onID,
		b.stdout,
		map[string]string{
			"APP_ID":           appID,
			"APP_NAME":         appName,
			"PGP_KEY":          b.pgpKey,
			"PGP_KEY_PASSWORD": b.pgpKeyPassword,
			"ARCHITECTURE":     b.architecture,
			"MSYS2PACKAGES":    strings.Join(b.packages, " "),
			"GOMAIN":           b.goMain,
			"GOFLAGS":          b.goFlags,
			"GOGENERATE":       b.goGenerate,
			"MSYS2INCLUDE":     b.include,
		},
		b.Render,
		[]string{},
	)
}
