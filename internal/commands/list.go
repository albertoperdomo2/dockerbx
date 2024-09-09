package commands

import (
	"context"
	"fmt"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"
)

func ListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List dockerbx containers",
		Run:   runList,
	}
}

func runList(cmd *cobra.Command, args []string) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		fmt.Printf("Error creating Docker client: %v\n", err)
		return
	}

	containers, err := cli.ContainerList(ctx, container.ListOptions{All: true})
	if err != nil {
		fmt.Printf("Error listing containers: %v\n", err)
		return
	}

	var dockerbxContainers []types.Container
	for _, container := range containers {
		if value, exists := container.Labels["owned_by"]; exists {
			if value == "dockerbx" {
				dockerbxContainers = append(dockerbxContainers, container)
			}
		}
	}

	if len(dockerbxContainers) > 0 {
		idWidth, nameWidth, commandWidth, stateWidth := 20, 20, 30, 10
		format := fmt.Sprintf("%%-%ds%%-%ds%%-%ds%%-%ds\n", idWidth, nameWidth, commandWidth, stateWidth)
		fmt.Printf(format, "CONTAINER_ID", "NAME", "COMMAND", "STATE")

		for _, dockerbxContainer := range dockerbxContainers {
			name := ""
			if len(dockerbxContainer.Names) > 0 {
				name = strings.TrimPrefix(dockerbxContainer.Names[0], "/")
			}

			id := truncateString(dockerbxContainer.ID, idWidth)
			name = truncateString(name, nameWidth)
			command := truncateString(dockerbxContainer.Command, commandWidth)
			state := truncateString(dockerbxContainer.State, stateWidth)

			fmt.Printf(format, id, name, command, state)
		}
	} else {
		fmt.Printf("No containers owned by \"dockerbx\" found.\n")
		return
	}
}

func truncateString(s string, maxLength int) string {
	if len(s) <= maxLength {
		return s
	}
	return s[:maxLength-5] + "...  "
}

// List containers owned by dockerbx
func ListDockerBxContainers(ctx context.Context, cli *client.Client) ([]types.Container, error) {
	var dockerbxContainers []types.Container

	containers, err := cli.ContainerList(ctx, container.ListOptions{All: true})
	if err != nil {
		return nil, err
	}

	for _, container := range containers {
		if value, exists := container.Labels["owned_by"]; exists {
			if value == "dockerbx" {
				dockerbxContainers = append(dockerbxContainers, container)
			}
		}
	}

	return dockerbxContainers, nil
}
