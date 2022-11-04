package tests

import (
	"context"

	"github.com/docker/docker/client"
	"github.com/pojntfx/hydrapp/hydrapp-builder/pkg/executors"
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
		onOutput,
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
	onID     func(id string)
	onOutput func(shortID string, color string, timestamp int64, message string)
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
		b.onOutput,
		map[string]string{
			"GOFLAGS":    b.goFlags,
			"GOGENERATE": b.goGenerate,
			"GOTESTS":    b.goTests,
		},
		b.Render,
		[]string{
			"sh",
			"-c",
			`export GOFLAGS="${GOFLAGS}" && cd /work && ${GOGENERATE} && cd /work && ${GOTESTS}`,
		},
	)
}
