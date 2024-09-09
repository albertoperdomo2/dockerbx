package main

import (
	"fmt"
	"os"

	"github.com/albertoperdomo2/dockerbx/internal/commands"
	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "dockerbx",
		Short: "dockerbx is a Docker-based alternative to toolbx",
		Long:  `A Docker-based tool for creating and managing containers for development environments.`,
	}

	rootCmd.AddCommand(commands.CreateCmd())
	rootCmd.AddCommand(commands.EnterCmd())
	rootCmd.AddCommand(commands.ListCmd())
	rootCmd.AddCommand(commands.RemoveCmd())
	rootCmd.AddCommand(commands.RunCmd())
	rootCmd.AddCommand(commands.UpdateCmd())
	rootCmd.AddCommand(commands.InitCmd())

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
