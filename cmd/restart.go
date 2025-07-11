package cmd

import (
	"fmt"

	"github.com/mohammed-bageri/dockenv/internal/config"
	"github.com/mohammed-bageri/dockenv/internal/docker"
	"github.com/mohammed-bageri/dockenv/internal/utils"

	"github.com/spf13/cobra"
)

var restartCmd = &cobra.Command{
	Use:   "restart [service...]",
	Short: "Restart development services",
	Long: `Restart the configured development services.
You can optionally specify which services to restart, otherwise all
configured services will be restarted.

Examples:
  dockenv restart           # Restart all services
  dockenv restart mysql     # Restart only MySQL
  dockenv restart mysql redis  # Restart MySQL and Redis`,
	RunE: runRestart,
}

func init() {
	rootCmd.AddCommand(restartCmd)
}

func runRestart(cmd *cobra.Command, args []string) error {
	// Check if Docker Compose file exists
	composePath := config.GetComposePath()
	if !utils.FileExists(composePath) {
		fmt.Println("❌ No Docker Compose file found.")
		fmt.Println("   Run 'dockenv init' first to set up your environment.")
		return fmt.Errorf("docker Compose file not found")
	}

	// Check Docker
	dockerInfo, err := docker.CheckDocker()
	if err != nil {
		return fmt.Errorf("failed to check Docker: %w", err)
	}

	if !dockerInfo.IsReady() {
		fmt.Println("❌ Docker setup incomplete:")
		fmt.Println(dockerInfo.GetInstallInstructions())
		return fmt.Errorf("docker setup required")
	}

	// Load config to validate services
	cfg, err := utils.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Validate requested services
	if len(args) > 0 {
		for _, serviceName := range args {
			if !utils.Contains(cfg.Services, serviceName) {
				return fmt.Errorf("service '%s' not configured. Available services: %v", serviceName, cfg.Services)
			}
		}
	}

	fmt.Println("🔄 Restarting services...")

	// Restart services
	if err := docker.ComposeRestart(composePath, args...); err != nil {
		return fmt.Errorf("failed to restart services: %w", err)
	}

	fmt.Println("✅ Services restarted successfully!")

	// Show what was restarted
	if len(args) == 0 {
		fmt.Printf("   Restarted: %v\n", cfg.Services)
	} else {
		fmt.Printf("   Restarted: %v\n", args)
	}

	fmt.Println("\nNext steps:")
	fmt.Println("  dockenv status   # Check service status")
	fmt.Println("  dockenv logs     # View service logs")

	return nil
}
