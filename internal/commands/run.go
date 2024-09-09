package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

func RunCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "enter [container_name]",
		Short: "Run a container",
		Run:   runRun,
	}
}

func runRun(cmd *cobra.Command, args []string) {

	fmt.Printf("Runing a container...")
}
