package commands

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"
)

func RemoveCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rm [container_name...]",
		Short: "Remove one or more containers",
		Run:   runRemove,
	}

	cmd.Flags().BoolP("force", "f", false, "Force removal of running containers")
	cmd.Flags().BoolP("all", "a", false, "Remove all dockerbx containers")

	return cmd
}

func runRemove(cmd *cobra.Command, args []string) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		fmt.Printf("Error creating Docker client: %v\n", err)
		return
	}

	force, _ := cmd.Flags().GetBool("force")
	all, _ := cmd.Flags().GetBool("all")

	var containerNames []string
	if all {
		containers, err := ListDockerBxContainers(ctx, cli)
		if err != nil {
			fmt.Printf("Error listing containers: %v\n", err)
			return
		}

		for _, c := range containers {
			containerNames = append(containerNames, c.ID)
		}
	} else if len(args) > 0 {
		containerNames = args
	} else {
		fmt.Printf("Error removing containers: No container name(s) provided.\n")
		return
	}

	for _, containerName := range containerNames {
		containerJSON, err := cli.ContainerInspect(ctx, containerName)
		if err != nil {
			fmt.Printf("Error inspecting container: %v\n", err)
			return
		}

		if containerJSON.State.Running && !force {
			fmt.Printf("Error removing container: Container is running. Use -f to force remove.\n")
			return
		}

		if containerJSON.State.Running && force {
			err = cli.ContainerStop(ctx, containerJSON.ID, container.StopOptions{})
			if err != nil {
				fmt.Printf("Error stoping container: %v\n", err)
				return
			}
		}

		err = cli.ContainerRemove(ctx, containerJSON.ID, container.RemoveOptions{Force: force})
		if err != nil {
			fmt.Printf("Error removing container: %v\n", err)
			return
		}

		fmt.Printf("Successfully removed container \"%v\"\n", containerName)
	}
}
