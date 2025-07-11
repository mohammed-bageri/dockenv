package cmd

import (
	"fmt"
	"strings"

	"github.com/mohammed-bageri/dockenv/internal/services"
	"github.com/mohammed-bageri/dockenv/internal/templates"
	"github.com/mohammed-bageri/dockenv/internal/utils"

	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add <service> [service...]",
	Short: "Add services to the configuration",
	Long: `Add one or more services to your existing dockenv configuration.
This will update the Docker Compose file and restart the environment.

Available services: ` + strings.Join(services.GetServiceNames(), ", ") + `

Examples:
  dockenv add mysql         # Add MySQL
  dockenv add redis mongodb # Add Redis and MongoDB
  dockenv add --port mysql:3307 mysql  # Add MySQL on custom port`,
	Args: cobra.MinimumNArgs(1),
	RunE: runAdd,
}

var addPortFlag []string

func init() {
	rootCmd.AddCommand(addCmd)

	addCmd.Flags().StringSliceVar(&addPortFlag, "port", []string{}, "Custom ports in format service:port")
}

func runAdd(cmd *cobra.Command, args []string) error {
	// Validate services
	if err := services.ValidateServices(args); err != nil {
		return err
	}

	// Load current configuration
	cfg, err := utils.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Check which services are new
	var newServices []string
	var existingServices []string

	for _, serviceName := range args {
		if utils.Contains(cfg.Services, serviceName) {
			existingServices = append(existingServices, serviceName)
		} else {
			newServices = append(newServices, serviceName)
		}
	}

	if len(existingServices) > 0 {
		fmt.Printf("⚠️  Already configured: %s\n", strings.Join(existingServices, ", "))
	}

	if len(newServices) == 0 {
		fmt.Println("✅ No new services to add.")
		return nil
	}

	fmt.Printf("➕ Adding services: %s\n", strings.Join(newServices, ", "))

	// Parse custom ports
	customPorts := make(map[string]int)
	for _, portSpec := range addPortFlag {
		parts := strings.Split(portSpec, ":")
		if len(parts) != 2 {
			return fmt.Errorf("invalid port specification: %s (expected format: service:port)", portSpec)
		}

		serviceName := parts[0]
		portStr := parts[1]

		port := 0
		if _, err := fmt.Sscanf(portStr, "%d", &port); err != nil {
			return fmt.Errorf("invalid port number: %s", portStr)
		}

		if !utils.Contains(newServices, serviceName) {
			return fmt.Errorf("service %s not in services to add", serviceName)
		}

		customPorts[serviceName] = port
	}

	// Add services to config
	for _, serviceName := range newServices {
		cfg.Services = append(cfg.Services, serviceName)

		// Set port
		if customPort, exists := customPorts[serviceName]; exists {
			cfg.Ports[serviceName] = customPort
		} else {
			service, _ := services.GetService(serviceName)
			cfg.Ports[serviceName] = service.DefaultPort
		}

		// Add environment variables
		service, _ := services.GetService(serviceName)
		for key, value := range service.EnvVars {
			if cfg.Env[key] == "" {
				cfg.Env[key] = value
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

	fmt.Println("✅ Services added successfully!")
	fmt.Printf("   Current services: %s\n", strings.Join(cfg.Services, ", "))

	// Show connection info for new services
	fmt.Println("\n📝 New Service Connection Information:")
	for _, serviceName := range newServices {
		port := cfg.Ports[serviceName]

		switch serviceName {
		case "mysql":
			fmt.Printf("  MySQL:      mysql://dockenv:password@localhost:%d/dockenv\n", port)
		case "postgres":
			fmt.Printf("  PostgreSQL: postgresql://dockenv:password@localhost:%d/dockenv\n", port)
		case "redis":
			fmt.Printf("  Redis:      redis://localhost:%d\n", port)
		case "mongodb":
			fmt.Printf("  MongoDB:    mongodb://dockenv:password@localhost:%d/dockenv\n", port)
		case "kafka":
			fmt.Printf("  Kafka:      localhost:%d\n", port)
		case "elasticsearch":
			fmt.Printf("  Elasticsearch: http://localhost:%d\n", port)
		case "rabbitmq":
			fmt.Printf("  RabbitMQ:   amqp://dockenv:password@localhost:%d\n", port)
			fmt.Printf("  RabbitMQ UI: http://localhost:15672 (admin panel)\n")
		}
	}

	fmt.Println("\nNext steps:")
	fmt.Println("  dockenv up       # Start all services (including new ones)")
	fmt.Println("  dockenv restart  # Restart to apply changes")

	return nil
}
