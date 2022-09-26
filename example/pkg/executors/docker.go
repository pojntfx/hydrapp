package executors

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

func DockerRunImage(
	ctx context.Context,
	cli *client.Client,
	image string,
	pull bool,
	privileged bool,
	dst string,
	env map[string]string,
) error {
	if pull {
		reader, err := cli.ImagePull(ctx, image, types.ImagePullOptions{})
		if err != nil {
			return err
		}

		if _, err := io.Copy(os.Stderr, reader); err != nil {
			return err
		}
	}

	pwd, err := os.Getwd()
	if err != nil {
		return err
	}

	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image:        image,
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
		OpenStdin:    true,
		Tty:          true,
		Env: func() []string {
			out := []string{}
			for key, value := range env {
				out = append(out, key+"="+value)
			}

			return out
		}(),
	}, &container.HostConfig{
		Privileged: privileged,
		Binds: []string{
			pwd + ":/work:z",
			dst + ":/dst:z",
		},
	}, nil, nil, "")
	if err != nil {
		return err
	}

	waiter, err := cli.ContainerAttach(ctx, resp.ID, types.ContainerAttachOptions{
		Stdin:  true,
		Stdout: true,
		Stderr: true,
		Stream: true,
	})
	if err != nil {
		return err
	}

	go io.Copy(waiter.Conn, os.Stdin)
	go io.Copy(os.Stdout, waiter.Reader)
	go io.Copy(os.Stderr, waiter.Reader)

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return err
	}

	statusChan, errChan := cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errChan:
		return err
	case status := <-statusChan:
		if status.Error != nil {
			return fmt.Errorf("could not wait for container: %v", status)
		}
	}

	return nil
}
