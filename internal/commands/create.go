package commands

import (
	"context"
	"fmt"
	"os"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"
)

func CreateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "create [container_name]",
		Short: "Create a new container",
		Run:   runCreate,
	}
}

func runCreate(cmd *cobra.Command, args []string) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		fmt.Printf("Error creating Docker client: %v\n", err)
		return
	}

	containerName := "dockerbx-default" // set this in config?
	if len(args) > 0 {
		containerName = args[0]
	}

	// Pull the latest Fedora image (this can be configured)
	_, err = cli.ImagePull(ctx, "fedora:latest", image.PullOptions{})
	if err != nil {
		fmt.Printf("Error pulling Fedora image: %v\n", err)
		return
	}

	// Get user's home dir
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("Error getting user's home directory: %v\n", err)
		return
	}

	fmt.Println("Creating container...")
	resp, err := cli.ContainerCreate(ctx,
		&container.Config{
			Image: "fedora:latest",
			Cmd:   []string{"/bin/bash"},
			Tty:   true,
			Labels: map[string]string{
				"owned_by": "dockerbx",
			},
		},
		&container.HostConfig{
			Mounts: []mount.Mount{
				{
					Type:   mount.TypeBind,
					Source: homeDir,
					Target: "/home/user",
				},
			},
		},
		nil,
		nil,
		containerName,
	)
	if err != nil {
		fmt.Printf("Error creating container: %v\n", err)
		return
	}

	fmt.Printf("Container created: %s\n", resp.ID)

	fmt.Println("Starting container...")
	if err := cli.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		fmt.Printf("Error starting container: %v\n", err)
		return
	}

	fmt.Printf("Container %s is running\n", containerName)
}
