package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/mohammed-bageri/dockenv/internal/config"
	"github.com/mohammed-bageri/dockenv/internal/docker"
	"github.com/mohammed-bageri/dockenv/internal/services"
	"github.com/mohammed-bageri/dockenv/internal/templates"
	"github.com/mohammed-bageri/dockenv/internal/utils"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new dockenv configuration",
	Long: `Interactive setup wizard that helps you choose services for your
development environment. This will create a configuration file and 
generate a Docker Compose file.

You can also use profiles for common setups:
  dockenv init --profile laravel  # MySQL + Redis
  dockenv init --profile node     # PostgreSQL + Redis
  dockenv init --profile django   # PostgreSQL + Redis
  dockenv init --profile full     # All services`,
	RunE: runInit,
}

var (
	profileFlag    string
	servicesFlag   []string
	autoDetectFlag bool
	portFlag       []string
	dataPathFlag   string
)

func init() {
	rootCmd.AddCommand(initCmd)

	initCmd.Flags().StringVar(&profileFlag, "profile", "", "Use a predefined service profile")
	initCmd.Flags().StringSliceVar(&servicesFlag, "services", []string{}, "Specify services directly")
	initCmd.Flags().BoolVar(&autoDetectFlag, "auto-detect", false, "Auto-detect project type and suggest services")
	initCmd.Flags().StringSliceVar(&portFlag, "port", []string{}, "Custom ports in format service:port")
	initCmd.Flags().StringVar(&dataPathFlag, "data-path", "", "Custom data directory path")
}

func runInit(cmd *cobra.Command, args []string) error {
	fmt.Println("🐳 Welcome to dockenv!")
	fmt.Println("   Setting up your local development environment...")
	fmt.Println()

	// Check Docker installation
	dockerInfo, err := docker.CheckDocker()
	if err != nil {
		return fmt.Errorf("failed to check Docker: %w", err)
	}

	if !dockerInfo.IsReady() {
		fmt.Println("❌ Docker setup incomplete:")
		fmt.Println(dockerInfo.GetInstallInstructions())
		fmt.Println()
		if !utils.PromptConfirm("Continue anyway?") {
			return fmt.Errorf("Docker setup required")
		}
	} else {
		fmt.Println("✅ Docker is ready!")
		fmt.Printf("   %s\n", dockerInfo.DockerVersion)
		fmt.Printf("   %s\n", dockerInfo.ComposeVersion)
		fmt.Println()
	}

	cfg, err := utils.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Set custom data path if provided
	if dataPathFlag != "" {
		cfg.DataPath = dataPathFlag
	}

	var selectedServices []string

	// Handle different initialization modes
	if profileFlag != "" {
		// Profile mode
		profileServices, exists := services.GetProfileServices(profileFlag)
		if !exists {
			return fmt.Errorf("unknown profile: %s. Available profiles: %s",
				profileFlag, strings.Join(services.GetProfileNames(), ", "))
		}
		selectedServices = profileServices
		fmt.Printf("📋 Using profile: %s\n", profileFlag)
		fmt.Printf("   Services: %s\n", strings.Join(selectedServices, ", "))

	} else if len(servicesFlag) > 0 {
		// Direct services mode
		if err := services.ValidateServices(servicesFlag); err != nil {
			return err
		}
		selectedServices = servicesFlag
		fmt.Printf("📋 Using specified services: %s\n", strings.Join(selectedServices, ", "))

	} else {
		// Interactive mode
		if autoDetectFlag {
			projectType := utils.DetectProjectType()
			if projectType != "unknown" {
				fmt.Printf("🔍 Detected project type: %s\n", projectType)
				if profileServices, exists := services.GetProfileServices(projectType); exists {
					if utils.PromptConfirm(fmt.Sprintf("Use recommended services for %s? (%s)",
						projectType, strings.Join(profileServices, ", "))) {
						selectedServices = profileServices
					}
				}
			}
		}

		if len(selectedServices) == 0 {
			// Interactive service selection
			selectedServices, err = selectServicesInteractive()
			if err != nil {
				return err
			}
		}
	}

	if len(selectedServices) == 0 {
		return fmt.Errorf("no services selected")
	}

	// Update configuration
	cfg.Services = selectedServices

	// Handle custom ports
	if err := parseCustomPorts(cfg); err != nil {
		return err
	}

	// Set default ports and env vars
	for _, serviceName := range selectedServices {
		service, _ := services.GetService(serviceName)

		if cfg.Ports[serviceName] == 0 {
			cfg.Ports[serviceName] = service.DefaultPort
		}

		for key, value := range service.EnvVars {
			if cfg.Env[key] == "" {
				cfg.Env[key] = value
			}
		}
	}

	// Create data directory
	if err := config.EnsureDataDir(); err != nil {
		return fmt.Errorf("failed to create data directory: %w", err)
	}

	// Save configuration
	if err := utils.SaveConfig(cfg); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	// Generate Docker Compose file
	if err := templates.GenerateDockerComposeEmbedded(cfg); err != nil {
		return fmt.Errorf("failed to generate Docker Compose file: %w", err)
	}

	// Generate .env file
	if err := utils.CreateEnvFile(cfg.Env); err != nil {
		return fmt.Errorf("failed to create .env file: %w", err)
	}

	fmt.Println()
	fmt.Println("✅ Configuration complete!")
	fmt.Printf("   Services: %s\n", strings.Join(cfg.Services, ", "))
	fmt.Printf("   Config: %s\n", config.GetConfigPath())
	fmt.Printf("   Compose: %s\n", config.GetComposePath())
	fmt.Printf("   Data: %s\n", cfg.DataPath)
	fmt.Printf("   Environment: %s\n", config.EnvFileName)
	fmt.Println()
	fmt.Println("Next steps:")
	fmt.Println("  dockenv up      # Start services")
	fmt.Println("  dockenv status  # Check service status")

	return nil
}

func selectServicesInteractive() ([]string, error) {
	fmt.Println("📋 Select services for your development environment:")
	fmt.Println()

	serviceOptions := []string{}
	serviceMap := make(map[string]string)

	for name, service := range services.AvailableServices {
		option := fmt.Sprintf("%s - %s (port %d)", service.DisplayName, service.Description, service.DefaultPort)
		serviceOptions = append(serviceOptions, option)
		serviceMap[option] = name
	}

	var selectedServices []string

	for {
		fmt.Println("\nSelect a service (or 'Done' to finish):")

		// Add "Done" option
		currentOptions := append(serviceOptions, "Done")

		selectPrompt := promptui.Select{
			Label: "Services",
			Items: currentOptions,
		}

		_, result, err := selectPrompt.Run()
		if err != nil {
			return nil, fmt.Errorf("selection failed: %w", err)
		}

		if result == "Done" {
			break
		}

		if serviceName, exists := serviceMap[result]; exists {
			if !utils.Contains(selectedServices, serviceName) {
				selectedServices = append(selectedServices, serviceName)
				fmt.Printf("✅ Added: %s\n", serviceName)
			} else {
				fmt.Printf("⚠️  Already selected: %s\n", serviceName)
			}
		}
	}

	return selectedServices, nil
}

func parseCustomPorts(cfg *config.Config) error {
	for _, portSpec := range portFlag {
		parts := strings.Split(portSpec, ":")
		if len(parts) != 2 {
			return fmt.Errorf("invalid port specification: %s (expected format: service:port)", portSpec)
		}

		serviceName := parts[0]
		portStr := parts[1]

		port, err := strconv.Atoi(portStr)
		if err != nil {
			return fmt.Errorf("invalid port number: %s", portStr)
		}

		if !utils.Contains(cfg.Services, serviceName) {
			return fmt.Errorf("service %s not in selected services", serviceName)
		}

		cfg.Ports[serviceName] = port
	}

	return nil
}
