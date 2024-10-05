package commands

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
	"github.com/spf13/cobra"
)

func BuildCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "build",
		Short: "Build a new image from a Dockerfile",
		Run:   runBuild,
	}

	cmd.Flags().String("file", "", "Path to the Dockerfile")
	cmd.Flags().String("name", "", "Built image name")
	cmd.MarkFlagRequired("file")
	cmd.MarkFlagRequired("name")

	return cmd
}

func runBuild(cmd *cobra.Command, args []string) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		fmt.Printf("Error creating Docker client: %v\n", err)
		return
	}

	dockerfile, _ := cmd.Flags().GetString("file")
	name, _ := cmd.Flags().GetString("name")

	build_ctx, err := archive.TarWithOptions(dockerfile, &archive.TarOptions{})
	if err != nil {
		fmt.Printf("Error creating tar context: %v\n", err)
		return
	}

	resp, err := cli.ImageBuild(ctx, build_ctx, types.ImageBuildOptions{
		Dockerfile: dockerfile,
		Tags:       []string{name},
		Remove:     true,
	})
	if err != nil {
		fmt.Printf("Error building image: %v\n", err)
		return
	}
	defer resp.Body.Close()

	fmt.Printf("Building image %s...\n", name)
	err = printBuildProgress(resp.Body)
	if err != nil {
		fmt.Printf("Error reading Docker build response: %v\n", err)
		return
	}

	fmt.Printf("\nYou can now run:\n")
	fmt.Printf("dockerbx create <container_name> --image %s\n", name)
}

func printBuildProgress(reader io.Reader) error {
	decoder := json.NewDecoder(reader)
	for {
		var message struct {
			Stream string `json:"stream"`
			Status string `json:"status"`
			Error  string `json:"error"`
		}

		if err := decoder.Decode(&message); err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}

		if message.Error != "" {
			return fmt.Errorf("build error: %s", message.Error)
		}

		if message.Stream != "" {
			fmt.Print(strings.TrimSpace(message.Stream))
		} else if message.Status != "" {
			fmt.Print(strings.TrimSpace(message.Status))
		}

		if message.Stream != "" || message.Status != "" {
			fmt.Println()
		}
	}
}
