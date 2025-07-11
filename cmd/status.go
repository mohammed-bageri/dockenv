package cmd

import (
	"fmt"

	"github.com/mohammed-bageri/dockenv/internal/config"
	"github.com/mohammed-bageri/dockenv/internal/docker"
	"github.com/mohammed-bageri/dockenv/internal/utils"

	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show service status",
	Long: `Display the current status of all configured services.
This shows which services are running, stopped, or have issues.

Examples:
  dockenv status       # Show status of all services`,
	RunE: runStatus,
}

func init() {
	rootCmd.AddCommand(statusCmd)
}

func runStatus(cmd *cobra.Command, args []string) error {
	// Check if Docker Compose file exists
	composePath := config.GetComposePath()
	if !utils.FileExists(composePath) {
		fmt.Println("❌ No Docker Compose file found.")
		fmt.Println("   Run 'dockenv init' first to set up your environment.")
		return nil
	}

	// Load configuration
	cfg, err := utils.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	if len(cfg.Services) == 0 {
		fmt.Println("📋 No services configured.")
		fmt.Println("   Run 'dockenv init' to set up services.")
		return nil
	}

	fmt.Println("📊 Service Status:")
	fmt.Printf("   Configuration: %s\n", config.GetConfigPath())
	fmt.Printf("   Compose file: %s\n", composePath)
	fmt.Printf("   Data directory: %s\n", cfg.DataPath)
	fmt.Println()

	// Check Docker
	dockerInfo, err := docker.CheckDocker()
	if err != nil {
		fmt.Printf("❌ Docker check failed: %v\n", err)
		return nil
	}

	if !dockerInfo.IsReady() {
		fmt.Println("❌ Docker not ready:")
		fmt.Println(dockerInfo.GetInstallInstructions())
		return nil
	}

	fmt.Println("🐳 Docker Status:")
	fmt.Printf("   %s\n", dockerInfo.DockerVersion)
	fmt.Printf("   %s\n", dockerInfo.ComposeVersion)
	fmt.Println()

	// Show configured services
	fmt.Println("🎯 Configured Services:")
	for _, serviceName := range cfg.Services {
		port := cfg.Ports[serviceName]
		fmt.Printf("   %-12s (port %d)\n", serviceName, port)
	}
	fmt.Println()

	// Show container status
	fmt.Println("📦 Container Status:")
	if err := docker.ComposeStatus(composePath); err != nil {
		fmt.Printf("   Failed to get status: %v\n", err)
		fmt.Println("   Services may not be running. Try 'dockenv up' to start them.")
	}

	return nil
}
