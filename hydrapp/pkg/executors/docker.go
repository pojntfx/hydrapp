package executors

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/user"
	"runtime"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	cp "github.com/otiai10/copy"
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
	stdout io.Writer,
	env map[string]string,
	renderTemplates func(workdir string, ejecting bool) error,
	cmds []string,
) error {
	currentUser, err := user.Current()
	if err != nil {
		return err
	}

	if runtime.GOOS != "windows" {
		env["DST_UID"] = currentUser.Uid
		env["DST_GID"] = currentUser.Gid
	}

	images, err := cli.ImageList(ctx, types.ImageListOptions{})
	if err != nil {
		return err
	}

	imageExists := false
o:
	for _, i := range images {
		for _, t := range i.RepoTags {
			if t == image {
				imageExists = true

				break o
			}
		}
	}

	if pull || !imageExists {
		reader, err := cli.ImagePull(ctx, image, types.ImagePullOptions{})
		if err != nil {
			return err
		}

		if _, err := io.Copy(stdout, reader); err != nil {
			return err
		}
	}

	workdir, err := os.MkdirTemp("", "hydrapp-build-dir-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(workdir)

	if err := cp.Copy(src, workdir); err != nil {
		return err
	}

	if err := renderTemplates(workdir, false); err != nil {
		return err
	}

	var cmd []string
	binds := []string{
		workdir + ":/work:z",
	}

	if len(cmds) > 0 {
		cmd = cmds
	} else {
		binds = append(binds, dst+":/hydrapp/dst:z")
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
		Cmd: cmd,
	}, &container.HostConfig{
		Privileged: privileged,
		Binds:      binds,
	}, nil, nil, "")
	if err != nil {
		return err
	}

	onID(resp.ID)

	waiter, err := cli.ContainerAttach(ctx, resp.ID, container.AttachOptions{
		Stdin:  false,
		Stdout: true,
		Stderr: true,
		Stream: true,
	})
	if err != nil {
		return err
	}
	defer waiter.Close()

	go io.Copy(stdout, waiter.Reader) // We intentionally ignore errors here since we can't handle them

	if err := cli.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
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
