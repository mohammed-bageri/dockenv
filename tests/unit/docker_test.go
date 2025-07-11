package unit

import (
	"strings"
	"testing"

	"github.com/mohammed-bageri/dockenv/internal/docker"
)

func TestCheckDocker(t *testing.T) {
	info, err := docker.CheckDocker()
	if err != nil {
		t.Errorf("CheckDocker() error = %v, want nil", err)
		return
	}

	if info == nil {
		t.Errorf("CheckDocker() returned nil info")
		return
	}

	// These tests depend on the environment, so we'll just check the structure
	t.Logf("Docker installed: %v", info.DockerInstalled)
	t.Logf("Docker Compose installed: %v", info.DockerComposeInstalled)
	t.Logf("Docker running: %v", info.DockerRunning)
	t.Logf("Docker version: %s", info.DockerVersion)
	t.Logf("Compose version: %s", info.ComposeVersion)

	// If docker is installed, version should not be empty
	if info.DockerInstalled && info.DockerVersion == "" {
		t.Errorf("Docker is installed but version is empty")
	}

	// If docker compose is installed, version should not be empty
	if info.DockerComposeInstalled && info.ComposeVersion == "" {
		t.Errorf("Docker Compose is installed but version is empty")
	}
}

func TestDockerInfo_IsReady(t *testing.T) {
	tests := []struct {
		name     string
		info     *docker.DockerInfo
		expected bool
	}{
		{
			name: "all ready",
			info: &docker.DockerInfo{
				DockerInstalled:        true,
				DockerComposeInstalled: true,
				DockerRunning:          true,
			},
			expected: true,
		},
		{
			name: "docker not installed",
			info: &docker.DockerInfo{
				DockerInstalled:        false,
				DockerComposeInstalled: true,
				DockerRunning:          true,
			},
			expected: false,
		},
		{
			name: "compose not installed",
			info: &docker.DockerInfo{
				DockerInstalled:        true,
				DockerComposeInstalled: false,
				DockerRunning:          true,
			},
			expected: false,
		},
		{
			name: "docker not running",
			info: &docker.DockerInfo{
				DockerInstalled:        true,
				DockerComposeInstalled: true,
				DockerRunning:          false,
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.info.IsReady()
			if result != tt.expected {
				t.Errorf("IsReady() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestDockerInfo_GetInstallInstructions(t *testing.T) {
	tests := []struct {
		name          string
		info          *docker.DockerInfo
		shouldContain []string
		shouldNotBe   string
	}{
		{
			name: "all ready",
			info: &docker.DockerInfo{
				DockerInstalled:        true,
				DockerComposeInstalled: true,
				DockerRunning:          true,
			},
			shouldContain: []string{"properly installed"},
		},
		{
			name: "docker not installed",
			info: &docker.DockerInfo{
				DockerInstalled:        false,
				DockerComposeInstalled: true,
				DockerRunning:          true,
			},
			shouldContain: []string{"Docker is not installed"},
		},
		{
			name: "docker not running",
			info: &docker.DockerInfo{
				DockerInstalled:        true,
				DockerComposeInstalled: true,
				DockerRunning:          false,
			},
			shouldContain: []string{"Docker daemon is not running"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			instructions := tt.info.GetInstallInstructions()

			for _, expected := range tt.shouldContain {
				if !strings.Contains(instructions, expected) {
					t.Errorf("GetInstallInstructions() should contain '%s', got: %s", expected, instructions)
				}
			}
		})
	}
}

func TestRunCompose(t *testing.T) {
	// We can't easily test this without a valid docker-compose file
	// But we can test that the function exists and doesn't panic with invalid args
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("RunCompose panicked: %v", r)
		}
	}()

	// This will likely fail, but should not panic
	_ = docker.RunCompose("--help")
}

func TestComposeValidate(t *testing.T) {
	// Test that the function doesn't panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("ComposeValidate panicked: %v", r)
		}
	}()

	// This will likely fail due to no file, but should not panic
	_ = docker.ComposeValidate("non-existent-file.yaml")
}
