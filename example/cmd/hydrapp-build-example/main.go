package main

import (
	"context"
	"flag"

	"github.com/docker/docker/client"
	"github.com/pojntfx/hydrapp/example/pkg/executors"
)

func main() {
	image := flag.String("image", "ghcr.io/pojntfx/hydrapp-build-example", "OCI image to use")
	pull := flag.Bool("pull", true, "Whether to pull the image or not")
	dst := flag.String("dst", "out", "Output directory")
	message := flag.String("message", "Hello, world!", "Message to print in the container")

	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}
	defer cli.Close()

	if err := executors.DockerRunImage(
		ctx,
		cli,
		*image,
		*pull,
		false,
		*dst,
		map[string]string{
			"MESSAGE": *message,
		},
	); err != nil {
		panic(err)
	}
}
