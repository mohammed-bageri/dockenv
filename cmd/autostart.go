package cmd

import (
	"fmt"

	"github.com/mohammed-bageri/dockenv/internal/systemd"

	"github.com/spf13/cobra"
)

var autostartCmd = &cobra.Command{
	Use:   "autostart",
	Short: "Manage auto-start configuration",
	Long: `Enable or disable automatic startup of services on system boot.
This creates a systemd service that will automatically start your
dockenv services when the system boots.

Note: This requires sudo privileges to create the systemd service.`,
}

var enableAutostartCmd = &cobra.Command{
	Use:   "enable",
	Short: "Enable auto-start on system boot",
	Long: `Enable automatic startup of dockenv services on system boot.
This creates a systemd service unit file and enables it.

Requirements:
- sudo privileges
- systemd-based system (most Linux distributions)
- Docker daemon configured to start on boot

The service will start services in the current directory, so make sure
to run this command from your project directory.`,
	RunE: runEnableAutostart,
}

var disableAutostartCmd = &cobra.Command{
	Use:   "disable",
	Short: "Disable auto-start on system boot",
	Long: `Disable automatic startup of dockenv services on system boot.
This stops and removes the systemd service unit file.

Requirements:
- sudo privileges`,
	RunE: runDisableAutostart,
}

var statusAutostartCmd = &cobra.Command{
	Use:   "status",
	Short: "Show auto-start status",
	Long: `Show the current status of the dockenv autostart service.
This will display whether the service is enabled, active, and show
recent log entries.`,
	RunE: runStatusAutostart,
}

func init() {
	rootCmd.AddCommand(autostartCmd)

	autostartCmd.AddCommand(enableAutostartCmd)
	autostartCmd.AddCommand(disableAutostartCmd)
	autostartCmd.AddCommand(statusAutostartCmd)
}

func runEnableAutostart(cmd *cobra.Command, args []string) error {
	fmt.Println("🔧 Enabling auto-start on system boot...")
	fmt.Println("   This requires sudo privileges.")
	fmt.Println()

	if err := systemd.EnableAutostart(); err != nil {
		return fmt.Errorf("failed to enable autostart: %w", err)
	}

	return nil
}

func runDisableAutostart(cmd *cobra.Command, args []string) error {
	fmt.Println("🔧 Disabling auto-start on system boot...")
	fmt.Println("   This requires sudo privileges.")
	fmt.Println()

	if err := systemd.DisableAutostart(); err != nil {
		return fmt.Errorf("failed to disable autostart: %w", err)
	}

	return nil
}

func runStatusAutostart(cmd *cobra.Command, args []string) error {
	fmt.Println("📊 Auto-start Status:")
	fmt.Println()

	enabled := systemd.IsEnabled()
	active := systemd.IsActive()

	if enabled {
		fmt.Println("✅ Auto-start is ENABLED")
	} else {
		fmt.Println("❌ Auto-start is DISABLED")
	}

	if active {
		fmt.Println("✅ Service is ACTIVE")
	} else {
		fmt.Println("❌ Service is INACTIVE")
	}

	fmt.Println()
	fmt.Println("📋 Detailed Status:")

	if err := systemd.GetStatus(); err != nil {
		fmt.Printf("Failed to get detailed status: %v\n", err)
	}

	fmt.Println()
	if !enabled {
		fmt.Println("To enable auto-start:")
		fmt.Println("  dockenv autostart enable")
	} else {
		fmt.Println("To disable auto-start:")
		fmt.Println("  dockenv autostart disable")
	}

	return nil
}
