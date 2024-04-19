package tests

import (
	"context"
	"fmt"
	"io"

	"github.com/docker/docker/client"
	"github.com/pojntfx/hydrapp/hydrapp/pkg/executors"
)

const (
	Image = "ghcr.io/pojntfx/hydrapp-build-tests"
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
	goFlags, // Flags to pass to the Go command
	goGenerate, // Command to execute go generate with
	goTests string, // Command to run tests with
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
		goFlags,
		goGenerate,
		goTests,
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
	goFlags,
	goGenerate,
	goTests string
}

func (b *Builder) Render(workdir string, ejecting bool) error {
	return nil
}

func (b *Builder) Build() error {
	return executors.DockerRunImage(
		b.ctx,
		b.cli,
		b.image,
		b.pull,
		true,
		b.src,
		"",
		b.onID,
		b.stdout,
		map[string]string{
			"GOFLAGS": b.goFlags,
		},
		b.Render,
		[]string{
			"sh",
			"-c",
			fmt.Sprintf(`export GOPROXY='https://proxy.golang.org,direct' GOFLAGS="${GOFLAGS}" && cd /work && %v && cd /work && %v`, b.goGenerate, b.goTests),
		},
	)
}
