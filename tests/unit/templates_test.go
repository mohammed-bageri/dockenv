package unit

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mohammed-bageri/dockenv/internal/config"
	"github.com/mohammed-bageri/dockenv/internal/templates"
	"github.com/mohammed-bageri/dockenv/internal/utils"
)

func TestGenerateDockerCompose(t *testing.T) {
	tempDir := t.TempDir()
	composeFile := filepath.Join(tempDir, "docker-compose.dockenv.yaml")

	// Create a test config
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
		DataPath: tempDir,
	}

	// Set environment variable to use temp directory
	oldCwd, _ := os.Getwd()
	defer os.Chdir(oldCwd)
	os.Chdir(tempDir)

	// Use embedded templates instead of external files
	err := templates.GenerateDockerComposeEmbedded(cfg)
	if err != nil {
		t.Errorf("GenerateDockerComposeEmbedded() error = %v, want nil", err)
		return
	}

	// Check if file was created
	if _, err := os.Stat(composeFile); os.IsNotExist(err) {
		t.Errorf("Docker compose file was not created: %s", composeFile)
		return
	}

	// Read and verify content
	content, err := os.ReadFile(composeFile)
	if err != nil {
		t.Errorf("Failed to read compose file: %v", err)
		return
	}

	contentStr := string(content)

	// Check for basic structure
	expectedStrings := []string{
		"version:",
		"services:",
		"mysql:",
		"redis:",
		"volumes:",
	}

	for _, expected := range expectedStrings {
		if !strings.Contains(contentStr, expected) {
			t.Errorf("Generated compose file should contain '%s'", expected)
		}
	}
}

func TestCreateEnvFile(t *testing.T) {
	tests := []struct {
		name            string
		initialContent  string
		newVars         map[string]string
		expectedContent string
		description     string
	}{
		{
			name:           "create_new_env_file",
			initialContent: "",
			newVars: map[string]string{
				"DB_HOST": "127.0.0.1",
				"DB_PORT": "3306",
			},
			expectedContent: "", // Will check individual variables since order may vary
			description:     "should create new .env file when none exists",
		},
		{
			name: "merge_with_existing_env",
			initialContent: `# Database configuration
DB_HOST=localhost
DB_PORT=3306
# Comment to preserve
EXISTING_VAR=value
`,
			newVars: map[string]string{
				"DB_HOST":    "127.0.0.1", // Should update existing
				"REDIS_HOST": "127.0.0.1", // Should add new
				"REDIS_PORT": "6379",      // Should add new
			},
			expectedContent: "", // Will check individual variables and comments
			description:     "should preserve existing content and comments while merging new variables",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := t.TempDir()
			envFile := filepath.Join(tempDir, ".env")

			// Set working directory to temp directory
			oldCwd, _ := os.Getwd()
			defer os.Chdir(oldCwd)
			os.Chdir(tempDir)

			// Create initial .env file if content provided
			if tt.initialContent != "" {
				err := os.WriteFile(envFile, []byte(tt.initialContent), 0644)
				if err != nil {
					t.Fatalf("Failed to create initial .env file: %v", err)
				}
			}

			// Call CreateEnvFile
			err := utils.CreateEnvFile(tt.newVars)
			if err != nil {
				t.Errorf("CreateEnvFile() error = %v, want nil", err)
				return
			}

			// Check if file was created
			if _, err := os.Stat(envFile); os.IsNotExist(err) {
				t.Errorf("Env file was not created: %s", envFile)
				return
			}

			// Read and verify content
			content, err := os.ReadFile(envFile)
			if err != nil {
				t.Errorf("Failed to read env file: %v", err)
				return
			}

			contentStr := string(content)

			// Test based on test case type
			if tt.name == "create_new_env_file" {
				// For new file creation, check that required variables are present
				requiredVars := []string{"DB_HOST=127.0.0.1", "DB_PORT=3306"}
				for _, required := range requiredVars {
					if !strings.Contains(contentStr, required) {
						t.Errorf("New env file should contain '%s'", required)
					}
				}
			} else if tt.name == "merge_with_existing_env" {
				// For merging test, check preserved comments and updated/new variables
				expectedElements := []string{
					"# Database configuration",
					"# Comment to preserve",
					"DB_HOST=127.0.0.1",    // Updated value
					"DB_PORT=3306",         // Preserved
					"EXISTING_VAR=value",   // Preserved
					"REDIS_HOST=127.0.0.1", // New
					"REDIS_PORT=6379",      // New
				}
				for _, expected := range expectedElements {
					if !strings.Contains(contentStr, expected) {
						t.Errorf("Merged env file should contain '%s'.\nActual content:\n%s", expected, contentStr)
					}
				}
			}
		})
	}
}

func TestGetEmbeddedTemplate(t *testing.T) {
	tests := []struct {
		name        string
		serviceName string
		expectError bool
	}{
		{"mysql template", "mysql", false},
		{"postgres template", "postgres", false},
		{"redis template", "redis", false},
		{"invalid service", "invalid_service", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			template, err := templates.GetEmbeddedTemplate(tt.serviceName)

			if tt.expectError {
				if err == nil {
					t.Errorf("GetEmbeddedTemplate(%v) expected error, got nil", tt.serviceName)
				}
				return
			}

			if err != nil {
				t.Errorf("GetEmbeddedTemplate(%v) error = %v, want nil", tt.serviceName, err)
				return
			}

			if template == "" {
				t.Errorf("GetEmbeddedTemplate(%v) returned empty template", tt.serviceName)
			}

			// Check that template contains expected service name
			if !strings.Contains(template, tt.serviceName) {
				t.Errorf("Template should contain service name '%s'", tt.serviceName)
			}
		})
	}
}

func TestGenerateDockerComposeWithInvalidService(t *testing.T) {
	tempDir := t.TempDir()

	// Create a config with invalid service
	cfg := &config.Config{
		Version:  "1.0",
		Services: []string{"invalid_service"},
		Ports:    map[string]int{},
		Env:      map[string]string{},
		Volumes:  map[string]string{},
		DataPath: tempDir,
	}

	// Set environment variable to use temp directory
	oldCwd, _ := os.Getwd()
	defer os.Chdir(oldCwd)
	os.Chdir(tempDir)

	err := templates.GenerateDockerCompose(cfg)
	if err == nil {
		t.Errorf("GenerateDockerCompose() expected error for invalid service, got nil")
	}
}
