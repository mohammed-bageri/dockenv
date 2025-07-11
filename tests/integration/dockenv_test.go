package integration

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

const binaryName = "dockenv"

func TestMain(m *testing.M) {
	// Build the binary before running tests
	build := exec.Command("go", "build", "-o", binaryName, "../../main.go")
	if err := build.Run(); err != nil {
		panic("Failed to build binary: " + err.Error())
	}

	// Run tests
	code := m.Run()

	// Cleanup
	os.Remove(binaryName)
	os.Exit(code)
}

func TestDockenvVersion(t *testing.T) {
	cmd := exec.Command("./"+binaryName, "--version")
	output, err := cmd.Output()
	if err != nil {
		t.Fatalf("Failed to run dockenv --version: %v", err)
	}

	outputStr := string(output)
	if !strings.Contains(outputStr, "dockenv version") {
		t.Errorf("Version output should contain 'dockenv version', got: %s", outputStr)
	}
}

func TestDockenvHelp(t *testing.T) {
	cmd := exec.Command("./"+binaryName, "--help")
	output, err := cmd.Output()
	if err != nil {
		t.Fatalf("Failed to run dockenv --help: %v", err)
	}

	outputStr := string(output)
	expectedCommands := []string{"init", "up", "down", "add", "remove", "status", "list"}

	for _, command := range expectedCommands {
		if !strings.Contains(outputStr, command) {
			t.Errorf("Help output should contain command '%s', got: %s", command, outputStr)
		}
	}
}

func TestDockenvList(t *testing.T) {
	cmd := exec.Command("./"+binaryName, "list")
	output, err := cmd.Output()
	if err != nil {
		t.Fatalf("Failed to run dockenv list: %v", err)
	}

	outputStr := string(output)
	expectedServices := []string{"mysql", "postgres", "redis", "mongodb"}

	for _, service := range expectedServices {
		if !strings.Contains(outputStr, service) {
			t.Errorf("List output should contain service '%s', got: %s", service, outputStr)
		}
	}
}

func TestDockenvInitWithProfile(t *testing.T) {
	tempDir := t.TempDir()
	oldDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
	defer os.Chdir(oldDir)

	// Change to temp directory
	err = os.Chdir(tempDir)
	if err != nil {
		t.Fatalf("Failed to change to temp directory: %v", err)
	}

	// Set environment variables to use temp directory
	os.Setenv("DOCKENV_CONFIG", filepath.Join(tempDir, "dockenv.yaml"))
	os.Setenv("DOCKENV_DATA", filepath.Join(tempDir, "data"))
	defer func() {
		os.Unsetenv("DOCKENV_CONFIG")
		os.Unsetenv("DOCKENV_DATA")
	}()

	// Test init with laravel profile (use --force-no-docker to skip Docker checks)
	cmd := exec.Command(filepath.Join(oldDir, binaryName), "init", "--profile", "laravel", "--services", "mysql,redis")
	cmd.Stdin = strings.NewReader("y\n") // Auto-answer yes to prompts
	output, err := cmd.CombinedOutput()

	// If Docker is not running, the command might fail, but we still want to test the basic functionality
	outputStr := string(output)
	if err != nil && !strings.Contains(outputStr, "Docker") {
		t.Fatalf("Failed to run dockenv init --profile laravel: %v\nOutput: %s", err, output)
	}

	// If Docker error, skip file checks
	if strings.Contains(outputStr, "Docker") && strings.Contains(outputStr, "not running") {
		t.Skip("Skipping init test - Docker not available in test environment")
		return
	}

	// Check if config file was created
	configPath := filepath.Join(tempDir, "dockenv.yaml")
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Errorf("Config file was not created: %s", configPath)
	}

	// Check if docker-compose file was created
	composePath := filepath.Join(tempDir, "docker-compose.dockenv.yaml")
	if _, err := os.Stat(composePath); os.IsNotExist(err) {
		t.Errorf("Docker compose file was not created: %s", composePath)
	}
}

func TestDockenvStatus(t *testing.T) {
	tempDir := t.TempDir()
	oldDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
	defer os.Chdir(oldDir)

	// Change to temp directory
	err = os.Chdir(tempDir)
	if err != nil {
		t.Fatalf("Failed to change to temp directory: %v", err)
	}

	// Test status without configuration
	cmd := exec.Command(filepath.Join(oldDir, binaryName), "status")
	output, _ := cmd.CombinedOutput()

	// Status should work even without configuration (might show no services)
	outputStr := string(output)
	if !strings.Contains(outputStr, "dockenv") {
		t.Errorf("Status output should contain 'dockenv', got: %s", outputStr)
	}
}

func TestDockenvAdd(t *testing.T) {
	tempDir := t.TempDir()
	oldDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
	defer os.Chdir(oldDir)

	// Change to temp directory
	err = os.Chdir(tempDir)
	if err != nil {
		t.Fatalf("Failed to change to temp directory: %v", err)
	}

	// Set environment variables to use temp directory
	os.Setenv("DOCKENV_CONFIG", filepath.Join(tempDir, "dockenv.yaml"))
	os.Setenv("DOCKENV_DATA", filepath.Join(tempDir, "data"))
	defer func() {
		os.Unsetenv("DOCKENV_CONFIG")
		os.Unsetenv("DOCKENV_DATA")
	}()

	// Initialize with basic setup first
	initCmd := exec.Command(filepath.Join(oldDir, binaryName), "init", "--profile", "node", "--services", "postgres,redis")
	initCmd.Stdin = strings.NewReader("y\n") // Auto-answer yes to prompts
	initOutput, err := initCmd.CombinedOutput()

	// Skip test if Docker is not available
	if err != nil && strings.Contains(string(initOutput), "Docker") {
		t.Skip("Skipping add test - Docker not available in test environment")
		return
	}

	if err != nil {
		t.Fatalf("Failed to initialize dockenv: %v\nOutput: %s", err, initOutput)
	}

	// Test adding a service
	addCmd := exec.Command(filepath.Join(oldDir, binaryName), "add", "mysql")
	output, err := addCmd.CombinedOutput()
	if err != nil {
		// If the error is due to Docker not running, skip the test
		if strings.Contains(string(output), "Docker") {
			t.Skip("Skipping add test - Docker not available")
			return
		}
		t.Fatalf("Failed to add service: %v\nOutput: %s", err, output)
	}

	// Check that service was added to config
	outputStr := string(output)
	if !strings.Contains(outputStr, "mysql") || !strings.Contains(outputStr, "added") {
		t.Errorf("Add output should indicate mysql was added, got: %s", outputStr)
	}
}
