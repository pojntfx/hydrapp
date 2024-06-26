package binaries

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/docker/docker/client"
	"github.com/pojntfx/hydrapp/hydrapp/pkg/builders"
	"github.com/pojntfx/hydrapp/hydrapp/pkg/executors"
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
	stdout io.Writer, // Writer to handle container output
	appID string, // App ID to use
	pgpKey []byte, // PGP key contents
	pgpKeyPassword, // password for the PGP key
	appName string, // App name
	branchID, // Branch ID
	branchName string, // Branch name
	branchTimestamp time.Time, // Branch timestamp
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
		stdout,
		appID,
		base64.StdEncoding.EncodeToString(pgpKey),
		pgpKeyPassword,
		appName,
		branchID,
		branchName,
		branchTimestamp,
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
	onID   func(id string)
	stdout io.Writer
	appID,
	pgpKey,
	pgpKeyPassword,
	appName,
	branchID,
	branchName string
	branchTimestamp time.Time
	goMain,
	goFlags,
	goGenerate,
	goExclude string
	hostPackages []string
}

func (b *Builder) Render(workdir string, ejecting bool) error {
	return nil
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
		b.stdout,
		map[string]string{
			"APP_ID":                   appID,
			"PGP_KEY":                  b.pgpKey,
			"PGP_KEY_PASSWORD":         b.pgpKeyPassword,
			"APP_NAME":                 appName,
			"GOMAIN":                   b.goMain,
			"GOFLAGS":                  b.goFlags,
			"GOGENERATE":               b.goGenerate,
			"GOEXCLUDE":                b.goExclude,
			"HOST_PACKAGES":            strings.Join(b.hostPackages, " "),
			"BRANCH_ID":                b.branchID,
			"BRANCH_TIMESTAMP_RFC3339": fmt.Sprintf("%v", b.branchTimestamp.Format(time.RFC3339)),
		},
		b.Render,
		[]string{},
	)
}
