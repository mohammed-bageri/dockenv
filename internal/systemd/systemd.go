package systemd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/mohammed-bageri/dockenv/internal/config"
)

const systemdTemplate = `[Unit]
Description=Dockenv Development Services
Requires=docker.service
After=docker.service
StartLimitIntervalSec=0

[Service]
Type=oneshot
RemainAfterExit=yes
WorkingDirectory=%s
ExecStart=/usr/local/bin/dockenv up
ExecStop=/usr/local/bin/dockenv down
TimeoutStartSec=0
User=%s

[Install]
WantedBy=multi-user.target
`

func EnableAutostart() error {
	// Get current working directory
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	// Get current user
	user := os.Getenv("USER")
	if user == "" {
		return fmt.Errorf("failed to get current user")
	}

	// Create systemd service file content
	serviceContent := fmt.Sprintf(systemdTemplate, cwd, user)

	// Write service file to temporary location
	tmpServiceFile := "/tmp/dockenv.service"
	if err := os.WriteFile(tmpServiceFile, []byte(serviceContent), 0644); err != nil {
		return fmt.Errorf("failed to write service file: %w", err)
	}

	// Copy service file to systemd directory (requires sudo)
	systemdDir := "/etc/systemd/system"
	serviceFile := filepath.Join(systemdDir, config.SystemdService)

	if err := exec.Command("sudo", "cp", tmpServiceFile, serviceFile).Run(); err != nil {
		return fmt.Errorf("failed to copy service file (try running with sudo): %w", err)
	}

	// Reload systemd
	if err := exec.Command("sudo", "systemctl", "daemon-reload").Run(); err != nil {
		return fmt.Errorf("failed to reload systemd: %w", err)
	}

	// Enable service
	if err := exec.Command("sudo", "systemctl", "enable", "dockenv.service").Run(); err != nil {
		return fmt.Errorf("failed to enable service: %w", err)
	}

	// Clean up temporary file
	os.Remove(tmpServiceFile)

	fmt.Println("✅ Autostart enabled successfully!")
	fmt.Println("   Services will start automatically on system boot.")
	fmt.Printf("   Service file: %s\n", serviceFile)
	fmt.Println("\nYou can manage the service with:")
	fmt.Println("   sudo systemctl start dockenv")
	fmt.Println("   sudo systemctl stop dockenv")
	fmt.Println("   sudo systemctl status dockenv")

	return nil
}

func DisableAutostart() error {
	// Stop service if running (ignore errors as service might not be running)
	_ = exec.Command("sudo", "systemctl", "stop", "dockenv.service").Run()

	// Disable service
	if err := exec.Command("sudo", "systemctl", "disable", "dockenv.service").Run(); err != nil {
		return fmt.Errorf("failed to disable service: %w", err)
	}

	// Remove service file
	serviceFile := filepath.Join("/etc/systemd/system", config.SystemdService)
	if err := exec.Command("sudo", "rm", "-f", serviceFile).Run(); err != nil {
		return fmt.Errorf("failed to remove service file: %w", err)
	}

	// Reload systemd
	if err := exec.Command("sudo", "systemctl", "daemon-reload").Run(); err != nil {
		return fmt.Errorf("failed to reload systemd: %w", err)
	}

	fmt.Println("✅ Autostart disabled successfully!")
	fmt.Printf("   Removed service file: %s\n", serviceFile)

	return nil
}

func GetStatus() error {
	return exec.Command("systemctl", "status", "dockenv.service").Run()
}

func IsEnabled() bool {
	err := exec.Command("systemctl", "is-enabled", "dockenv.service").Run()
	return err == nil
}

func IsActive() bool {
	err := exec.Command("systemctl", "is-active", "dockenv.service").Run()
	return err == nil
}
