package commands

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/albertoperdomo2/dockerbx/internal/config"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

func EnterCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "enter [container_name]",
		Short: "Enter an existing container",
		Run:   runEnter,
	}
}

func runEnter(cmd *cobra.Command, args []string) {
	ctx := context.Background()

	config, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		return
	}

	containerName := config.DefaultName
	if len(args) > 0 {
		containerName = args[0]
	}

	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		fmt.Printf("Error creating Docker client: %v\n", err)
		return
	}

	_, err = cli.ContainerInspect(ctx, containerName)
	if err != nil {
		fmt.Printf("Container %s does not exist. Please create it first.\n", containerName)
		return
	}

	containerJSON, err := cli.ContainerInspect(ctx, containerName)
	if err != nil {
		fmt.Printf("Error inspecting container: %v\n", err)
		return
	}

	if !containerJSON.State.Running {
		fmt.Printf("Container %s is not running. Starting it now...\n", containerName)
		err = cli.ContainerStart(ctx, containerName, container.StartOptions{})
		if err != nil {
			fmt.Printf("Error starting container: %v\n", err)
			return
		}
	}

	// exec default config
	execConfig := types.ExecConfig{
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
		Tty:          true,
		Cmd:          []string{"/bin/sh"}, // change this
	}

	execID, err := cli.ContainerExecCreate(ctx, containerName, execConfig)
	if err != nil {
		fmt.Printf("Error creating exec instance: %v\n", err)
		return
	}

	resp, err := cli.ContainerExecAttach(ctx, execID.ID, types.ExecStartCheck{Tty: true})
	if err != nil {
		fmt.Printf("Error attaching to exec instance: %v\n", err)
		return
	}
	defer resp.Close()

	// setup terminal
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		fmt.Printf("Error setting up terminal: %v\n", err)
		return
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	err = cli.ContainerExecStart(ctx, execID.ID, types.ExecStartCheck{Tty: true})
	if err != nil {
		fmt.Printf("Error starting exec instance: %v\n", err)
		return
	}

	// handle I/O
	go func() {
		io.Copy(os.Stdout, resp.Reader)
	}()
	go func() {
		io.Copy(resp.Conn, os.Stdin)
	}()

	// wait for exec instance to finish
	for {
		inspectResp, err := cli.ContainerExecInspect(ctx, execID.ID)
		if err != nil {
			fmt.Printf("Error inspecting exec instance: %v\n", err)
			return
		}

		if !inspectResp.Running {
			break
		}
	}
}
