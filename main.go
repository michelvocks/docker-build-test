package main

import (
	"context"
	"io"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
)

func main() {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}
	cli.NegotiateAPIVersion(ctx)

	reader, err := cli.ImagePull(ctx, "golang:latest", types.ImagePullOptions{})
	if err != nil {
		panic(err)
	}
	io.Copy(os.Stdout, reader)

	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: "golang:latest",
		Cmd:   []string{"echo", "Hello World", "&&", "echo", "to another World"},
	}, &container.HostConfig{
		Mounts: []mount.Mount{
			mount.Mount{
				Type:   mount.TypeBind,
				Source: "/Users/michelvocks/go/src/github.com/michelvocks/docker-build-test/gaia-docker-test",
				Target: "/tmp",
			}},
	}, nil, "")
	if err != nil {
		panic(err)
	}

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}

	execID, err := cli.ContainerExecCreate(ctx, resp.ID, types.ExecConfig{
		/*Cmd: []string{
			"go",
			"get",
			"-d",
			"./...",
		},*/
		Cmd:        []string{"askjfhksdhjfkjshf"},
		WorkingDir: "/tmp",
	})
	if err != nil {
		panic(err)
	}

	if err := cli.ContainerExecStart(ctx, execID.ID, types.ExecStartCheck{}); err != nil {
		panic(err)
	}

	out, err := cli.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{ShowStdout: true})
	if err != nil {
		panic(err)
	}

	stdcopy.StdCopy(os.Stdout, os.Stderr, out)

	execID, err = cli.ContainerExecCreate(ctx, resp.ID, types.ExecConfig{
		Cmd: []string{
			"go",
			"run",
			"main.go",
		},
		WorkingDir: "/tmp",
	})
	if err != nil {
		panic(err)
	}

	if err := cli.ContainerExecStart(ctx, execID.ID, types.ExecStartCheck{}); err != nil {
		panic(err)
	}
}
