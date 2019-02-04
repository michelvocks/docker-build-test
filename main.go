package main

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
)

func main() {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.WithHost("unix:///var/run/docker.sock"))
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
		Cmd:   []string{"sleep", "1d"},
	}, &container.HostConfig{Mounts: []mount.Mount{
		mount.Mount{Source: "/tmp/mount", Target: "/tmp/mount"},
	}}, nil, "")
	if err != nil {
		panic(err)
	}

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}

	execID, err := cli.ContainerExecCreate(ctx, resp.ID, types.ExecConfig{
		Cmd:          []string{"git", "clone", "https://github.com/michelvocks/gaia-docker-test", "src"},
		AttachStdout: true,
		AttachStderr: true,
		WorkingDir:   "/tmp",
	})
	if err != nil {
		panic(err)
	}

	hijackResponse, err := cli.ContainerExecAttach(ctx, execID.ID, types.ExecStartCheck{})
	if err != nil {
		panic(err)
	}

	if err := cli.ContainerExecStart(ctx, execID.ID, types.ExecStartCheck{}); err != nil {
		panic(err)
	}

	b, err := ioutil.ReadAll(hijackResponse.Reader)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Git Clone output: %s\n", b)

	execID, err = cli.ContainerExecCreate(ctx, resp.ID, types.ExecConfig{
		Cmd: []string{
			"go",
			"get",
			"-d",
			"./...",
		},
		AttachStdout: true,
		AttachStderr: true,
		WorkingDir:   "/tmp/src",
	})
	if err != nil {
		panic(err)
	}

	hijackResponse, err = cli.ContainerExecAttach(ctx, execID.ID, types.ExecStartCheck{})
	if err != nil {
		panic(err)
	}

	if err := cli.ContainerExecStart(ctx, execID.ID, types.ExecStartCheck{}); err != nil {
		panic(err)
	}

	b, err = ioutil.ReadAll(hijackResponse.Reader)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Go get output: %s\n", b)

	execID, err = cli.ContainerExecCreate(ctx, resp.ID, types.ExecConfig{
		Cmd: []string{
			"go",
			"build",
			"-o",
			"/tmp/pipeline.app",
		},
		WorkingDir: "/tmp/src",
	})
	if err != nil {
		panic(err)
	}

	hijackResponse, err = cli.ContainerExecAttach(ctx, execID.ID, types.ExecStartCheck{})
	if err != nil {
		panic(err)
	}

	if err := cli.ContainerExecStart(ctx, execID.ID, types.ExecStartCheck{}); err != nil {
		panic(err)
	}

	b, err = ioutil.ReadAll(hijackResponse.Reader)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Go build output: %s\n", b)

	if err := cli.ContainerStop(ctx, resp.ID, nil); err != nil {
		panic(err)
	}

	commitResp, err := cli.ContainerCommit(ctx, resp.ID, types.ContainerCommitOptions{Reference: "helloworld"})
	if err != nil {
		panic(err)
	}
	fmt.Println(commitResp.ID)
}
