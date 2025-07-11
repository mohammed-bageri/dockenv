package unit

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/mohammed-bageri/dockenv/internal/config"
	"github.com/mohammed-bageri/dockenv/internal/utils"
)

func TestFileExists(t *testing.T) {
	// Create a temporary file
	tempFile, err := os.CreateTemp("", "test_file")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())
	tempFile.Close()

	tests := []struct {
		name     string
		filename string
		expected bool
	}{
		{"existing file", tempFile.Name(), true},
		{"non-existing file", "/path/that/does/not/exist", false},
		{"empty string", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := utils.FileExists(tt.filename)
			if result != tt.expected {
				t.Errorf("FileExists(%v) = %v, want %v", tt.filename, result, tt.expected)
			}
		})
	}
}

func TestLoadConfig(t *testing.T) {
	tests := []struct {
		name           string
		configContent  string
		setupConfig    bool
		expectError    bool
		expectedFields map[string]interface{}
	}{
		{
			name:        "no config file",
			setupConfig: false,
			expectError: false,
			expectedFields: map[string]interface{}{
				"version":  "1.0",
				"services": 0, // length of services slice
			},
		},
		{
			name: "valid config file",
			configContent: `version: "1.0"
services:
  - mysql
  - redis
ports:
  mysql: 3306
  redis: 6379
env:
  DB_HOST: "127.0.0.1"
data_path: "/custom/path"`,
			setupConfig: true,
			expectError: false,
			expectedFields: map[string]interface{}{
				"version":  "1.0",
				"services": 2,
			},
		},
		{
			name: "invalid yaml",
			configContent: `version: "1.0"
services:
  - mysql
    invalid: yaml structure`,
			setupConfig: true,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := t.TempDir()
			configPath := filepath.Join(tempDir, "dockenv.yaml")

			// Set environment variable to use our temp config
			os.Setenv("DOCKENV_CONFIG", configPath)
			defer os.Unsetenv("DOCKENV_CONFIG")

			if tt.setupConfig {
				err := os.WriteFile(configPath, []byte(tt.configContent), 0644)
				if err != nil {
					t.Fatalf("Failed to write test config: %v", err)
				}
			}

			cfg, err := utils.LoadConfig()

			if tt.expectError {
				if err == nil {
					t.Errorf("LoadConfig() expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("LoadConfig() unexpected error: %v", err)
				return
			}

			if cfg == nil {
				t.Errorf("LoadConfig() returned nil config")
				return
			}

			// Check expected fields
			if version, ok := tt.expectedFields["version"]; ok {
				if cfg.Version != version.(string) {
					t.Errorf("LoadConfig() version = %v, want %v", cfg.Version, version)
				}
			}

			if servicesLen, ok := tt.expectedFields["services"]; ok {
				if len(cfg.Services) != servicesLen.(int) {
					t.Errorf("LoadConfig() services length = %v, want %v", len(cfg.Services), servicesLen)
				}
			}

			// Ensure maps are initialized
			if cfg.Ports == nil {
				t.Errorf("LoadConfig() Ports map should be initialized")
			}
			if cfg.Env == nil {
				t.Errorf("LoadConfig() Env map should be initialized")
			}
			if cfg.Volumes == nil {
				t.Errorf("LoadConfig() Volumes map should be initialized")
			}
		})
	}
}

func TestSaveConfig(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "dockenv.yaml")

	// Set environment variable to use our temp config
	os.Setenv("DOCKENV_CONFIG", configPath)
	defer os.Unsetenv("DOCKENV_CONFIG")

	cfg := &config.Config{
		Version:  "1.0",
		Services: []string{"mysql", "redis"},
		Ports: map[string]int{
			"mysql": 3306,
			"redis": 6379,
		},
		Env: map[string]string{
			"DB_HOST": "127.0.0.1",
		},
		Volumes: map[string]string{
			"mysql_data": "/data/mysql",
		},
		DataPath: "/custom/path",
	}

	err := utils.SaveConfig(cfg)
	if err != nil {
		t.Errorf("SaveConfig() error = %v, want nil", err)
		return
	}

	// Verify file was created
	if !utils.FileExists(configPath) {
		t.Errorf("SaveConfig() did not create config file")
		return
	}

	// Load and verify content
	loadedCfg, err := utils.LoadConfig()
	if err != nil {
		t.Errorf("Failed to load saved config: %v", err)
		return
	}

	if loadedCfg.Version != cfg.Version {
		t.Errorf("Saved config version = %v, want %v", loadedCfg.Version, cfg.Version)
	}

	if len(loadedCfg.Services) != len(cfg.Services) {
		t.Errorf("Saved config services length = %v, want %v", len(loadedCfg.Services), len(cfg.Services))
	}
}
