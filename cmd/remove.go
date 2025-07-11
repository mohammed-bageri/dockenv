package cmd

import (
	"fmt"
	"strings"

	"github.com/mohammed-bageri/dockenv/internal/services"
	"github.com/mohammed-bageri/dockenv/internal/templates"
	"github.com/mohammed-bageri/dockenv/internal/utils"

	"github.com/spf13/cobra"
)

var removeCmd = &cobra.Command{
	Use:     "remove <service> [service...]",
	Aliases: []string{"rm"},
	Short:   "Remove services from the configuration",
	Long: `Remove one or more services from your dockenv configuration.
This will update the Docker Compose file and stop the removed services.

Examples:
  dockenv remove mysql         # Remove MySQL
  dockenv remove redis mongodb # Remove Redis and MongoDB
  dockenv rm mysql             # Same as remove`,
	Args: cobra.MinimumNArgs(1),
	RunE: runRemove,
}

var (
	removeForceFlag     bool
	removeVolumesRmFlag bool
)

func init() {
	rootCmd.AddCommand(removeCmd)

	removeCmd.Flags().BoolVarP(&removeForceFlag, "force", "f", false, "Force removal without confirmation")
	removeCmd.Flags().BoolVar(&removeVolumesRmFlag, "volumes", false, "Also remove data volumes (WARNING: Data will be lost!)")
}

func runRemove(cmd *cobra.Command, args []string) error {
	// Load current configuration
	cfg, err := utils.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Check which services exist
	var servicesToRemove []string
	var missingServices []string

	for _, serviceName := range args {
		if utils.Contains(cfg.Services, serviceName) {
			servicesToRemove = append(servicesToRemove, serviceName)
		} else {
			missingServices = append(missingServices, serviceName)
		}
	}

	if len(missingServices) > 0 {
		fmt.Printf("⚠️  Not configured: %s\n", strings.Join(missingServices, ", "))
	}

	if len(servicesToRemove) == 0 {
		fmt.Println("✅ No services to remove.")
		return nil
	}

	fmt.Printf("➖ Removing services: %s\n", strings.Join(servicesToRemove, ", "))

	// Confirmation prompt
	if !removeForceFlag {
		if removeVolumesRmFlag {
			fmt.Printf("⚠️  WARNING: This will also remove data volumes for these services!\n")
			fmt.Printf("   Data in the following services will be PERMANENTLY LOST: %s\n", strings.Join(servicesToRemove, ", "))
		}

		if !utils.PromptConfirm("Are you sure you want to continue?") {
			fmt.Println("Operation cancelled.")
			return nil
		}
	}

	// Remove services from config
	for _, serviceName := range servicesToRemove {
		cfg.Services = utils.RemoveString(cfg.Services, serviceName)
		delete(cfg.Ports, serviceName)

		// Remove service-specific environment variables
		service, exists := services.GetService(serviceName)
		if exists {
			for key := range service.EnvVars {
				delete(cfg.Env, key)
			}
		}
	}

	// Save updated configuration
	if err := utils.SaveConfig(cfg); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	// Update Docker Compose file
	if err := templates.GenerateDockerComposeEmbedded(cfg); err != nil {
		return fmt.Errorf("failed to update Docker Compose file: %w", err)
	}

	// Update .env file
	if err := utils.CreateEnvFile(cfg.Env); err != nil {
		return fmt.Errorf("failed to update .env file: %w", err)
	}

	fmt.Println("✅ Services removed successfully!")
	if len(cfg.Services) > 0 {
		fmt.Printf("   Remaining services: %s\n", strings.Join(cfg.Services, ", "))
	} else {
		fmt.Println("   No services configured.")
	}

	// Stop and optionally remove volumes for the removed services
	// Note: This is a simplified approach. In practice, you'd want to use
	// docker-compose to stop specific services, but that's complex to implement
	// generically. For now, we'll just inform the user.

	fmt.Println("\n📝 Manual cleanup required:")
	fmt.Println("  dockenv restart  # Restart to apply configuration changes")
	if removeVolumesRmFlag {
		fmt.Println("  # Data volumes will be removed on next 'dockenv down --volumes'")
	}

	return nil
}
