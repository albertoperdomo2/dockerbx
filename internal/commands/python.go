package commands

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/albertoperdomo2/dockerbx/internal/config"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"
)

func PythonCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "python [container_name]",
		Short: "Create a new Python environment",
		Run:   runPythonCreate,
	}

	cmd.Flags().String("version", "3.9", "Python version to use")
	cmd.Flags().String("venv", "", "Name of the virtual environment to create")
	cmd.Flags().StringSlice("packages", []string{}, "List of packages to install")
	cmd.Flags().String("requirements", "", "Path to requirements.txt file")

	return cmd
}

func runPythonCreate(cmd *cobra.Command, args []string) {
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

	pythonVersion, _ := cmd.Flags().GetString("version")
	venvName, _ := cmd.Flags().GetString("venv")
	packages, _ := cmd.Flags().GetStringSlice("packages")
	requirementsFile, _ := cmd.Flags().GetString("requirements")

	// in this case, use a Python base image
	baseImage := fmt.Sprintf("python:%s-slim", pythonVersion)

	reader, err := cli.ImagePull(ctx, baseImage, image.PullOptions{})
	if err != nil {
		fmt.Printf("Error pulling image: %v\n", err)
		return
	}
	defer reader.Close()

	// wait for the image pull to complete
	_, err = io.Copy(os.Stdout, reader)
	if err != nil {
		fmt.Printf("Error reading image pull response: %v\n", err)
		return
	}

	mounts := config.Mounts
	if requirementsFile != "" {
		absPath, err := filepath.Abs(requirementsFile)
		if err != nil {
			fmt.Printf("Error getting absolute path of requirements file: %v\n", err)
			return
		}
		mounts = append(mounts, mount.Mount{
			Type:   mount.TypeBind,
			Source: absPath,
			Target: "/tmp/requirements.txt",
		})
	}

	resp, err := cli.ContainerCreate(ctx,
		&container.Config{
			Image: baseImage,
			Cmd:   []string{"/bin/bash"},
			Tty:   true,
			Labels: map[string]string{
				"owned_by": "dockerbx",
				"type":     "python",
			},
			Env: []string{"PS1=\\[\\e[32m\\]⬢\\[\\e[0m\\][\\u@dockerbx-python](\\W)\\$ "},
		},
		&container.HostConfig{
			Mounts: mounts,
		},
		nil,
		nil,
		containerName,
	)
	if err != nil {
		fmt.Printf("Error creating container: %v\n", err)
		return
	}

	fmt.Printf("Python container created: %s\n", resp.ID)

	if err := cli.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		fmt.Printf("Error starting container: %v\n", err)
		return
	}

	if venvName != "" {
		createVenvCmd := fmt.Sprintf("python -m venv %s", venvName)
		execResp, err := cli.ContainerExecCreate(ctx, resp.ID, types.ExecConfig{
			Cmd: []string{"/bin/sh", "-c", createVenvCmd},
		})
		if err != nil {
			fmt.Printf("Error creating virtual environment: %v\n", err)
			return
		}
		if err := cli.ContainerExecStart(ctx, execResp.ID, types.ExecStartCheck{}); err != nil {
			fmt.Printf("Error starting virtual environment creation: %v\n", err)
			return
		}
		fmt.Printf("Virtual environment '%s' created\n", venvName)
	}

	if len(packages) > 0 || requirementsFile != "" {
		var installCmd string
		if len(packages) > 0 {
			installCmd = fmt.Sprintf("pip install %s", strings.Join(packages, " "))
		}
		if requirementsFile != "" {
			if installCmd != "" {
				installCmd += " && "
			}
			installCmd += "pip install -r /tmp/requirements.txt"
		}
		execResp, err := cli.ContainerExecCreate(ctx, resp.ID, types.ExecConfig{
			Cmd: []string{"/bin/sh", "-c", installCmd},
		})
		if err != nil {
			fmt.Printf("Error creating package installation command: %v\n", err)
			return
		}
		if err := cli.ContainerExecStart(ctx, execResp.ID, types.ExecStartCheck{}); err != nil {
			fmt.Printf("Error starting package installation: %v\n", err)
			return
		}
		if len(packages) > 0 {
			fmt.Printf("Packages installed: %s\n", strings.Join(packages, ", "))
		}
		if requirementsFile != "" {
			fmt.Printf("Packages installed from requirements.txt\n")
		}
	}

	fmt.Printf("Python container %s is running\n", containerName)

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
