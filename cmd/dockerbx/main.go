package main

import (
	"fmt"
	"os"
	"runtime/debug"

	"github.com/albertoperdomo2/dockerbx/internal/commands"
	"github.com/spf13/cobra"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
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
	rootCmd.AddCommand(versionCmd())

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func versionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the version number of dockerbx",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("dockerbx version %s\n", version)
			if info, ok := debug.ReadBuildInfo(); ok {
				fmt.Printf("go version %s\n", info.GoVersion)
				for _, setting := range info.Settings {
					switch setting.Key {
					case "vcs.revision":
						fmt.Printf("git commit: %s\n", setting.Value)
					case "vcs.time":
						fmt.Printf("build date: %s\n", setting.Value)
					}
				}
			}
		},
	}
}
