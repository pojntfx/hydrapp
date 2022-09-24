package main

import (
	"context"
	"flag"
	"io"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

func main() {
	image := flag.String("image", "ghcr.io/pojntfx/hydrapp-build-example", "OCI image to use")
	pull := flag.Bool("pull", true, "Whether to pull the image or not")
	message := flag.String("message", "Hello, world!", "Message to print in the container")

	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}

	if *pull {
		reader, err := cli.ImagePull(ctx, *image, types.ImagePullOptions{})
		if err != nil {
			panic(err)
		}

		if _, err := io.Copy(os.Stderr, reader); err != nil {
			panic(err)
		}
	}

	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image:        *image,
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
		OpenStdin:    true,
		Tty:          true,
		Env:          []string{"MESSAGE=" + *message},
	}, nil, nil, nil, "")
	if err != nil {
		panic(err)
	}

	waiter, err := cli.ContainerAttach(ctx, resp.ID, types.ContainerAttachOptions{
		Stdin:  true,
		Stdout: true,
		Stderr: true,
		Stream: true,
	})
	if err != nil {
		panic(err)
	}

	go io.Copy(waiter.Conn, os.Stdin)
	go io.Copy(os.Stdout, waiter.Reader)
	go io.Copy(os.Stderr, waiter.Reader)

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}

	statusChan, errChan := cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errChan:
		panic(err)
	case status := <-statusChan:
		if status.Error != nil {
			panic(status)
		}
	}
}
