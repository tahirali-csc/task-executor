package dockerrunner

import (
	"bufio"
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/task-executor/pkg/docker-runner/api"
	"io"
	"io/ioutil"
	"log"
)

func Run(config api.ContainerConfig) error {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}

	ctx := context.Background()
	reader, err := cli.ImagePull(ctx, config.Image, types.ImagePullOptions{})
	if err != nil {
		return err
	}

	io.Copy(ioutil.Discard, reader)
	reader.Close()

	var volumeMounts []mount.Mount
	for _, v := range config.Volumes {
		volumeMounts = append(volumeMounts, mount.Mount{
			Target: v.Target,
			Source: v.Source,
			Type:   mount.TypeBind,
		})
	}

	resp, err := cli.ContainerCreate(context.Background(), &container.Config{
		Image:        config.Image,
		Cmd:          config.Command,
		Tty:          true,
		AttachStdout: true,
		AttachStderr: true,
		Env:          config.Env,
	}, &container.HostConfig{
		Mounts:     volumeMounts,
		AutoRemove: true,
	}, nil, "")
	if err != nil {
		return err
	}

	err = cli.ContainerStart(context.Background(), resp.ID, types.ContainerStartOptions{})
	if err != nil {
		return err
	}

	out, err := cli.ContainerLogs(context.Background(), resp.ID, types.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Follow:     true,
	})
	if err != nil {
		return err
	}

	rdr := bufio.NewReader(out)
	for {
		line, _, err := rdr.ReadLine()
		if err != nil {
			break
		}
		log.Print(string(line))
	}

	return nil
}
