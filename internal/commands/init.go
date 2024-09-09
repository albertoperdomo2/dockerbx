package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

func InitCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "enter [container_name]",
		Short: "Init a container",
		Run:   runInit,
	}
}

func runInit(cmd *cobra.Command, args []string) {

	fmt.Printf("Initing a container...")
}
