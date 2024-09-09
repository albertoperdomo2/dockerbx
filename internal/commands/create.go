package commands

import (
	"context"
	"fmt"

	"github.com/albertoperdomo2/dockerbx/internal/config"
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

	config, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		return
	}

	containerName := config.DefaultName
	if len(args) > 0 {
		containerName = args[0]
	}

	_, err = cli.ImagePull(ctx, config.BaseImage, image.PullOptions{})
	if err != nil {
		fmt.Printf("Error pulling Fedora image: %v\n", err)
		return
	}

	if err != nil {
		fmt.Printf("Error getting user's home directory: %v\n", err)
		return
	}

	resp, err := cli.ContainerCreate(ctx,
		&container.Config{
			Image: config.BaseImage,
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
					Source: config.Mounts[0].Source,
					Target: config.Mounts[0].Target,
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

	if err := cli.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		fmt.Printf("Error starting container: %v\n", err)
		return
	}

	fmt.Printf("Container %s is running\n", containerName)
}
