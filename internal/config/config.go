package config

import (
	"fmt"
	"os"
	"path/filepath"
)

const (
	ConfigFileName  = "dockenv.yaml"
	ComposeFileName = "docker-compose.dockenv.yaml"
	EnvFileName     = ".env"
	SystemdService  = "dockenv.service"
	DefaultDataPath = "/var/lib/dockenv"
)

type Config struct {
	Version  string            `yaml:"version"`
	Services []string          `yaml:"services"`
	Ports    map[string]int    `yaml:"ports,omitempty"`
	Env      map[string]string `yaml:"env,omitempty"`
	Volumes  map[string]string `yaml:"volumes,omitempty"`
	DataPath string            `yaml:"data_path,omitempty"`
}

func GetConfigPath() string {
	if configPath := os.Getenv("DOCKENV_CONFIG"); configPath != "" {
		return configPath
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "./" + ConfigFileName
	}

	return filepath.Join(homeDir, ".config", "dockenv", ConfigFileName)
}

func GetComposePath() string {
	return "./" + ComposeFileName
}

func GetDataPath() string {
	if dataPath := os.Getenv("DOCKENV_DATA"); dataPath != "" {
		return dataPath
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return DefaultDataPath
	}

	return filepath.Join(homeDir, ".local", "share", "dockenv")
}

func EnsureConfigDir() error {
	configPath := GetConfigPath()
	configDir := filepath.Dir(configPath)

	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	return nil
}

func EnsureDataDir() error {
	dataPath := GetDataPath()

	if err := os.MkdirAll(dataPath, 0755); err != nil {
		return fmt.Errorf("failed to create data directory: %w", err)
	}

	return nil
}
