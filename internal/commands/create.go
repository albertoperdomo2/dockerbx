package commands

import (
	"context"
	"fmt"

	"github.com/albertoperdomo2/dockerbx/internal/config"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"
)

func CreateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create [container_name]",
		Short: "Create a new container",
		Run:   runCreate,
	}

	cmd.Flags().String("clone", "", "Git repository URL to clone")

	return cmd
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
			Env: []string{"PS1=\\[\\e[32m\\]⬢\\[\\e[0m\\][\\u@dockerbx](\\W)\\$ "},
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

	cloneURL, _ := cmd.Flags().GetString("clone")
	if cloneURL != "" {
		cloneCmd := fmt.Sprintf("git clone %s /app", cloneURL)
		execResp, err := cli.ContainerExecCreate(ctx, resp.ID, types.ExecConfig{
			Cmd: []string{"/bin/sh", "-c", cloneCmd},
		})
		if err != nil {
			fmt.Printf("Error creating clone command: %v\n", err)
			return
		}
		if err := cli.ContainerExecStart(ctx, execResp.ID, types.ExecStartCheck{}); err != nil {
			fmt.Printf("Error starting clone: %v\n", err)
			return
		}
		fmt.Printf("Repository cloned: %s\n", cloneURL)
	}
}
