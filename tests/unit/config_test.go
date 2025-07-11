package unit

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mohammed-bageri/dockenv/internal/config"
)

func TestGetConfigPath(t *testing.T) {
	tests := []struct {
		name        string
		envVar      string
		expectedEnd string
	}{
		{
			name:        "default path",
			envVar:      "",
			expectedEnd: ".config/dockenv/dockenv.yaml",
		},
		{
			name:        "custom env path",
			envVar:      "/custom/path/dockenv.yaml",
			expectedEnd: "/custom/path/dockenv.yaml",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.envVar != "" {
				os.Setenv("DOCKENV_CONFIG", tt.envVar)
				defer os.Unsetenv("DOCKENV_CONFIG")
			}

			result := config.GetConfigPath()
			if tt.envVar != "" {
				if result != tt.expectedEnd {
					t.Errorf("GetConfigPath() = %v, want %v", result, tt.expectedEnd)
				}
			} else {
				if !filepath.IsAbs(result) {
					t.Errorf("GetConfigPath() should return absolute path, got %v", result)
				}
				if !strings.HasSuffix(result, tt.expectedEnd) {
					t.Errorf("GetConfigPath() = %v, should end with %v", result, tt.expectedEnd)
				}
			}
		})
	}
}

func TestGetComposePath(t *testing.T) {
	expected := "./docker-compose.dockenv.yaml"
	result := config.GetComposePath()
	if result != expected {
		t.Errorf("GetComposePath() = %v, want %v", result, expected)
	}
}

func TestGetDataPath(t *testing.T) {
	tests := []struct {
		name        string
		envVar      string
		expectedEnd string
	}{
		{
			name:        "default path",
			envVar:      "",
			expectedEnd: ".local/share/dockenv",
		},
		{
			name:        "custom env path",
			envVar:      "/custom/data/path",
			expectedEnd: "/custom/data/path",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.envVar != "" {
				os.Setenv("DOCKENV_DATA", tt.envVar)
				defer os.Unsetenv("DOCKENV_DATA")
			}

			result := config.GetDataPath()
			if tt.envVar != "" {
				if result != tt.expectedEnd {
					t.Errorf("GetDataPath() = %v, want %v", result, tt.expectedEnd)
				}
			} else {
				if !filepath.IsAbs(result) {
					t.Errorf("GetDataPath() should return absolute path, got %v", result)
				}
				if !strings.HasSuffix(result, tt.expectedEnd) {
					t.Errorf("GetDataPath() = %v, should end with %v", result, tt.expectedEnd)
				}
			}
		})
	}
}

func TestEnsureConfigDir(t *testing.T) {
	tempDir := t.TempDir()
	os.Setenv("DOCKENV_CONFIG", filepath.Join(tempDir, "config", "dockenv.yaml"))
	defer os.Unsetenv("DOCKENV_CONFIG")

	err := config.EnsureConfigDir()
	if err != nil {
		t.Errorf("EnsureConfigDir() error = %v, want nil", err)
	}

	configDir := filepath.Dir(config.GetConfigPath())
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		t.Errorf("Config directory was not created: %v", configDir)
	}
}

func TestEnsureDataDir(t *testing.T) {
	tempDir := t.TempDir()
	os.Setenv("DOCKENV_DATA", filepath.Join(tempDir, "data"))
	defer os.Unsetenv("DOCKENV_DATA")

	err := config.EnsureDataDir()
	if err != nil {
		t.Errorf("EnsureDataDir() error = %v, want nil", err)
	}

	dataDir := config.GetDataPath()
	if _, err := os.Stat(dataDir); os.IsNotExist(err) {
		t.Errorf("Data directory was not created: %v", dataDir)
	}
}
