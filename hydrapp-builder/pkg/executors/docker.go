package executors

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"
	"time"
	"unicode"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/pojntfx/hydrapp/hydrapp-builder/pkg/utils"
)

func DockerRunImage(
	ctx context.Context,
	cli *client.Client,
	image string,
	pull bool,
	privileged bool,
	src string,
	dst string,
	onID func(id string),
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
			src + ":/src:ro",
			dst + ":/dst:z",
		},
	}, nil, nil, "")
	if err != nil {
		return err
	}

	onID(resp.ID)

	waiter, err := cli.ContainerAttach(ctx, resp.ID, types.ContainerAttachOptions{
		Stdin:  false,
		Stdout: true,
		Stderr: true,
		Stream: true,
	})
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(waiter.Reader)
	go func() {
		c := utils.GetRandomANSIColor()

		for scanner.Scan() {
			if runtime.GOOS == "windows" {
				fmt.Printf(
					"%v@%v %v\n",
					resp.ID[:4],
					time.Now().Unix(),
					strings.TrimFunc(scanner.Text(), func(r rune) bool {
						return !unicode.IsGraphic(r)
					}),
				)
			} else {
				fmt.Printf(
					"%v%v%v@%v%v %v%v%v\n",
					utils.ColorBackgroundBlack,
					c,
					resp.ID[:4],
					time.Now().Unix(),
					utils.ColorReset,
					c,
					strings.TrimFunc(scanner.Text(), func(r rune) bool {
						return !unicode.IsGraphic(r)
					}),
					utils.ColorReset,
				)
			}
		}

		if scanner.Err() != nil {
			panic(err)
		}

		fmt.Printf("%v", utils.ColorReset)
	}()

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return err
	}

	statusChan, errChan := cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errChan:
		return err
	case status := <-statusChan:
		if (status.StatusCode != 0 && status.StatusCode != 137) || status.Error != nil {
			return fmt.Errorf("could not wait for container: %v", status)
		}
	}

	return nil
}
