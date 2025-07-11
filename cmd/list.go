package cmd

import (
	"fmt"
	"strings"

	"github.com/mohammed-bageri/dockenv/internal/config"
	"github.com/mohammed-bageri/dockenv/internal/services"
	"github.com/mohammed-bageri/dockenv/internal/utils"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List available services and profiles",
	Long: `Display all available services that can be configured with dockenv,
along with predefined profiles for common development stacks.

This command shows:
- Available services with their default ports
- Predefined profiles (laravel, node, django, etc.)
- Currently configured services (if any)`,
	RunE: runList,
}

var (
	listServicesFlag bool
	listProfilesFlag bool
	listCurrentFlag  bool
)

func init() {
	rootCmd.AddCommand(listCmd)

	listCmd.Flags().BoolVar(&listServicesFlag, "services", false, "Show only available services")
	listCmd.Flags().BoolVar(&listProfilesFlag, "profiles", false, "Show only available profiles")
	listCmd.Flags().BoolVar(&listCurrentFlag, "current", false, "Show only currently configured services")
}

func runList(cmd *cobra.Command, args []string) error {
	// Determine what to show based on flags
	showServices := listServicesFlag || (!listProfilesFlag && !listCurrentFlag)
	showProfiles := listProfilesFlag || (!listServicesFlag && !listCurrentFlag)
	showCurrent := listCurrentFlag || (!listServicesFlag && !listProfilesFlag)

	if showCurrent {
		if err := showCurrentServices(); err != nil {
			return err
		}
		fmt.Println()
	}

	if showServices {
		showAvailableServices()
		fmt.Println()
	}

	if showProfiles {
		showAvailableProfiles()
	}

	return nil
}

func showCurrentServices() error {
	fmt.Println("📋 Currently Configured Services:")

	cfg, err := utils.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	if len(cfg.Services) == 0 {
		fmt.Println("   No services configured. Run 'dockenv init' to get started.")
		return nil
	}

	fmt.Printf("   Configuration: %s\n", config.GetConfigPath())
	fmt.Printf("   Data path: %s\n", cfg.DataPath)
	fmt.Println()

	for _, serviceName := range cfg.Services {
		service, exists := services.GetService(serviceName)
		if !exists {
			fmt.Printf("   %-12s (unknown service)\n", serviceName)
			continue
		}

		port := cfg.Ports[serviceName]
		if port == 0 {
			port = service.DefaultPort
		}

		fmt.Printf("   %-12s %-20s port %d\n",
			service.DisplayName,
			service.Description,
			port)
	}

	return nil
}

func showAvailableServices() {
	fmt.Println("🛠️  Available Services:")

	for _, serviceName := range services.GetServiceNames() {
		service, _ := services.GetService(serviceName)
		fmt.Printf("   %-12s %-30s (default port %d)\n",
			serviceName,
			service.Description,
			service.DefaultPort)
	}

	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("   dockenv init --services mysql,redis")
	fmt.Println("   dockenv add postgres")
	fmt.Println("   dockenv remove mongodb")
}

func showAvailableProfiles() {
	fmt.Println("📦 Available Profiles:")

	profiles := services.Profiles
	for profileName, profileServices := range profiles {
		fmt.Printf("   %-10s %s\n", profileName, strings.Join(profileServices, ", "))
	}

	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("   dockenv init --profile laravel")
	fmt.Println("   dockenv init --profile node")
}
