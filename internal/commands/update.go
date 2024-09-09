package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

func UpdateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "enter [container_name]",
		Short: "Update a container",
		Run:   runUpdate,
	}
}

func runUpdate(cmd *cobra.Command, args []string) {

	fmt.Printf("Updateing a container...")
}
