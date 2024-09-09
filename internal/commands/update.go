package commands

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/albertoperdomo2/dockerbx/internal/config"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"
)

func UpdateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update [container_name]",
		Short: "Update a container's base image and packages",
		Run:   runUpdate,
	}

	cmd.Flags().BoolP("packages", "p", false, "Update packages within the container")

	return cmd
}

func runUpdate(cmd *cobra.Command, args []string) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		fmt.Printf("Error creating Docker client: %v\n", err)
		return
	}

	config, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		return
	}

	containerName := config.DefaultName
	if len(args) > 0 {
		containerName = args[0]
	}

	containerJSON, err := cli.ContainerInspect(ctx, containerName)
	if err != nil {
		fmt.Printf("Error: Container '%s' not found.\n", containerName)
		return
	}

	imageName := containerJSON.Config.Image

	fmt.Printf("Pulling latest version of %s...\n", imageName)
	reader, err := cli.ImagePull(ctx, imageName, image.PullOptions{})
	if err != nil {
		fmt.Printf("Error pulling latest image: %v\n", err)
		return
	}
	io.Copy(os.Stdout, reader)

	fmt.Println("Creating new container with updated image...")
	newContainerName := containerName + "-updated"
	newContainer, err := cli.ContainerCreate(ctx, containerJSON.Config, containerJSON.HostConfig, nil, nil, newContainerName)
	if err != nil {
		fmt.Printf("Error creating new container: %v\n", err)
		return
	}

	if containerJSON.State.Running {
		fmt.Println("Stopping the old container...")
		if err := cli.ContainerStop(ctx, containerName, container.StopOptions{}); err != nil {
			fmt.Printf("Error stopping old container: %v\n", err)
			return
		}
	}

	fmt.Println("Removing the old container...")
	if err := cli.ContainerRemove(ctx, containerName, container.RemoveOptions{}); err != nil {
		fmt.Printf("Error removing old container: %v\n", err)
		return
	}

	fmt.Println("Renaming the new container...")
	if err := cli.ContainerRename(ctx, newContainer.ID, containerName); err != nil {
		fmt.Printf("Error renaming new container: %v\n", err)
		return
	}

	updatePackages, _ := cmd.Flags().GetBool("packages")
	if updatePackages {
		fmt.Println("Updating packages within the container...")
		execConfig := types.ExecConfig{
			Cmd:          []string{"dnf", "update", "-y"},
			AttachStdout: true,
			AttachStderr: true,
		}

		execID, err := cli.ContainerExecCreate(ctx, containerName, execConfig)
		if err != nil {
			fmt.Printf("Error creating exec instance for package update: %v\n", err)
			return
		}

		resp, err := cli.ContainerExecAttach(ctx, execID.ID, types.ExecStartCheck{})
		if err != nil {
			fmt.Printf("Error starting exec instance for package update: %v\n", err)
			return
		}
		defer resp.Close()

		io.Copy(os.Stdout, resp.Reader)
	}

	fmt.Printf("Container '%s' has been updated successfully.\n", containerName)
}
