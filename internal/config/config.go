package config

import (
	"os"
	"path/filepath"

	"github.com/docker/docker/api/types/mount"
	"gopkg.in/yaml.v2"
)

type Config struct {
	BaseImage   string        `yaml:"base_image"`
	DefaultName string        `yaml:"default_name"`
	Mounts      []mount.Mount `yaml:"mounts"`
}

func LoadConfig() (*Config, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	configPath := filepath.Join(homeDir, ".config", "dockerbx", "dockerbx.yaml")
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
