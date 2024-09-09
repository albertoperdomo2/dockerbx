package commands

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"
)

func RunCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                "run [container_name] [-t] [command]",
		Short:              "Run a command in a container",
		DisableFlagParsing: true,
		Run:                runRun,
	}

	return cmd
}

func runRun(cmd *cobra.Command, args []string) {
	if len(args) < 2 {
		fmt.Println("Error: Please provide both a container name and a command to run.")
		return
	}

	containerName := args[0]
	command := args[1:]

	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		fmt.Printf("Error creating Docker client: %v\n", err)
		return
	}

	containerJSON, err := cli.ContainerInspect(ctx, containerName)
	if err != nil {
		fmt.Printf("Error: Container '%s' not found.\n", containerName)
		return
	}

	if !containerJSON.State.Running {
		fmt.Printf("Container '%s' is not running. Starting it now...\n", containerName)
		err = cli.ContainerStart(ctx, containerName, container.StartOptions{})
		if err != nil {
			fmt.Printf("Error starting container: %v\n", err)
			return
		}
	}

	tty := false
	if command[0] == "--tty" || command[0] == "-t" {
		tty = true
		command = command[1:]
	}

	execConfig := types.ExecConfig{
		Cmd:          command,
		AttachStdout: true,
		AttachStderr: true,
		Tty:          tty,
	}

	execID, err := cli.ContainerExecCreate(ctx, containerName, execConfig)
	if err != nil {
		fmt.Printf("Error creating exec instance: %v\n", err)
		return
	}

	resp, err := cli.ContainerExecAttach(ctx, execID.ID, types.ExecStartCheck{Tty: tty})
	if err != nil {
		fmt.Printf("Error attaching to exec instance: %v\n", err)
		return
	}
	defer resp.Close()

	if tty {
		_, err = io.Copy(os.Stdout, resp.Reader)
	} else {
		_, err = io.Copy(os.Stdout, resp.Reader)
	}
	if err != nil && err != io.EOF {
		fmt.Printf("Error streaming command output: %v\n", err)
	}

	inspectResp, err := cli.ContainerExecInspect(ctx, execID.ID)
	if err != nil {
		fmt.Printf("Error inspecting exec instance: %v\n", err)
		return
	}

	if inspectResp.ExitCode != 0 {
		fmt.Printf("Command exited with non-zero status: %d\n", inspectResp.ExitCode)
		os.Exit(inspectResp.ExitCode)
	}
}
