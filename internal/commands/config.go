package commands

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/albertoperdomo2/dockerbx/internal/config"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

func ExportConfigCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "export-config [file_name]",
		Short: "Export the current configuration",
		Run:   runExportConfig,
	}
}

func ImportConfigCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "import-config [file_name]",
		Short: "Import a configuration",
		Run:   runImportConfig,
	}
}

func runExportConfig(cmd *cobra.Command, args []string) {
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		return
	}

	fileName := "dockerbx_config.yaml"
	if len(args) > 0 {
		fileName = args[0]
	}

	data, err := yaml.Marshal(cfg)
	if err != nil {
		fmt.Printf("Error marshaling config: %v\n", err)
		return
	}

	err = ioutil.WriteFile(fileName, data, 0644)
	if err != nil {
		fmt.Printf("Error writing config file: %v\n", err)
		return
	}

	fmt.Printf("Configuration exported to %s\n", fileName)
}

func runImportConfig(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		fmt.Println("Please provide a file name to import")
		return
	}

	fileName := args[0]
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		fmt.Printf("Error reading config file: %v\n", err)
		return
	}

	var cfg config.Config
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		fmt.Printf("Error unmarshaling config: %v\n", err)
		return
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("Error getting user's home directory: %v\n", err)
		return
	}

	configPath := filepath.Join(homeDir, ".config", "dockerbx", "dockerbx.yaml")
	err = ioutil.WriteFile(configPath, data, 0644)
	if err != nil {
		fmt.Printf("Error writing config file: %v\n", err)
		return
	}

	fmt.Printf("Configuration imported from %s\n", fileName)
}
