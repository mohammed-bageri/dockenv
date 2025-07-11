package cmd

import (
	"fmt"

	"github.com/mohammed-bageri/dockenv/internal/config"
	"github.com/mohammed-bageri/dockenv/internal/docker"
	"github.com/mohammed-bageri/dockenv/internal/utils"

	"github.com/spf13/cobra"
)

var downCmd = &cobra.Command{
	Use:   "down",
	Short: "Stop development services",
	Long: `Stop all development services and remove containers.
Data volumes are preserved for persistent storage.

Examples:
  dockenv down              # Stop all services
  dockenv down --volumes    # Stop services and remove volumes (data loss!)`,
	RunE: runDown,
}

var (
	removeVolumesFlag bool
	removeImagesFlag  bool
)

func init() {
	rootCmd.AddCommand(downCmd)

	downCmd.Flags().BoolVarP(&removeVolumesFlag, "volumes", "v", false, "Remove volumes (WARNING: This will delete all data!)")
	downCmd.Flags().BoolVar(&removeImagesFlag, "rmi", false, "Remove images")
}

func runDown(cmd *cobra.Command, args []string) error {
	// Check if Docker Compose file exists
	composePath := config.GetComposePath()
	if !utils.FileExists(composePath) {
		fmt.Println("❌ No Docker Compose file found.")
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

	// Warning for volume removal
	if removeVolumesFlag {
		fmt.Println("⚠️  WARNING: This will permanently delete all data in volumes!")
		if !utils.PromptConfirm("Are you sure you want to continue?") {
			fmt.Println("Operation cancelled.")
			return nil
		}
	}

	fmt.Println("🛑 Stopping services...")

	// Stop services
	if err := docker.ComposeDown(composePath); err != nil {
		return fmt.Errorf("failed to stop services: %w", err)
	}

	// Handle additional cleanup
	if removeVolumesFlag {
		fmt.Println("🗑️  Removing volumes...")
		if err := docker.RunCompose("-f", composePath, "down", "-v"); err != nil {
			return fmt.Errorf("failed to remove volumes: %w", err)
		}
	}

	if removeImagesFlag {
		fmt.Println("🗑️  Removing images...")
		if err := docker.RunCompose("-f", composePath, "down", "--rmi", "all"); err != nil {
			return fmt.Errorf("failed to remove images: %w", err)
		}
	}

	fmt.Println("✅ Services stopped successfully!")

	if removeVolumesFlag {
		fmt.Println("   📁 All data volumes have been removed.")
	} else {
		fmt.Println("   📁 Data volumes preserved for next startup.")
	}

	fmt.Println("\nNext steps:")
	fmt.Println("  dockenv up       # Start services again")
	fmt.Println("  dockenv status   # Check service status")

	return nil
}
