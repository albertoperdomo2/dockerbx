package commands

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"
)

const (
	baseImage       = "fedora:latest"
	configFileName  = "dockerbx.yaml"
	dockerbxNetwork = "dockerbx-network"
)

func InitCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "init",
		Short: "Initialize dockerbx environment",
		Long:  `Set up necessary Docker images and configurations for dockerbx to function properly.`,
		Run:   runInit,
	}
}

func runInit(cmd *cobra.Command, args []string) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		fmt.Printf("Error creating Docker client: %v\n", err)
		return
	}

	// Pull base image
	fmt.Printf("Pulling base image %s...\n", baseImage)
	reader, err := cli.ImagePull(ctx, baseImage, image.PullOptions{})
	if err != nil {
		fmt.Printf("Error pulling base image: %v\n", err)
		return
	}
	io.Copy(os.Stdout, reader)

	// Create default configuration file
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("Error getting user home directory: %v\n", err)
		return
	}

	configDir := filepath.Join(homeDir, ".config", "dockerbx")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		fmt.Printf("Error creating config directory: %v\n", err)
		return
	}

	configPath := filepath.Join(configDir, configFileName)
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		fmt.Println("Creating default configuration file...")
		defaultConfig := []byte(`base_image: fedora:latest
default_name: dockerbx-default
mounts:
  - source: $HOME
    target: /home/user
`)
		if err := os.WriteFile(configPath, defaultConfig, 0644); err != nil {
			fmt.Printf("Error writing default configuration file: %v\n", err)
			return
		}
	}

	// Create Docker network
	fmt.Println("Creating dockerbx network...")
	_, err = cli.NetworkCreate(ctx, dockerbxNetwork, types.NetworkCreate{
		Driver: "bridge",
		Labels: map[string]string{"com.dockerbx.network": "true"},
	})
	if err != nil {
		fmt.Printf("Error creating dockerbx network: %v\n", err)
		return
	}

	fmt.Println("dockerbx initialized successfully!")
	fmt.Printf("Configuration file created at: %s\n", configPath)
	fmt.Println("You can now start using dockerbx to create and manage containers.")
}
