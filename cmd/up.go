package cmd

import (
	"fmt"

	"github.com/mohammed-bageri/dockenv/internal/config"
	"github.com/mohammed-bageri/dockenv/internal/docker"
	"github.com/mohammed-bageri/dockenv/internal/utils"

	"github.com/spf13/cobra"
)

var upCmd = &cobra.Command{
	Use:   "up [service...]",
	Short: "Start development services",
	Long: `Start the configured development services using Docker Compose.
You can optionally specify which services to start, otherwise all
configured services will be started.

Examples:
  dockenv up           # Start all services
  dockenv up mysql     # Start only MySQL
  dockenv up mysql redis  # Start MySQL and Redis`,
	RunE: runUp,
}

var (
	detachFlag        bool
	buildFlag         bool
	removeOrphansFlag bool
)

func init() {
	rootCmd.AddCommand(upCmd)

	upCmd.Flags().BoolVarP(&detachFlag, "detach", "d", true, "Run containers in detached mode")
	upCmd.Flags().BoolVar(&buildFlag, "build", false, "Build images before starting")
	upCmd.Flags().BoolVar(&removeOrphansFlag, "remove-orphans", false, "Remove containers for services not defined in compose file")
}

func runUp(cmd *cobra.Command, args []string) error {
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

	// Load config to get service info
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

	// Create data directories
	if err := config.EnsureDataDir(); err != nil {
		return fmt.Errorf("failed to ensure data directory: %w", err)
	}

	fmt.Println("🚀 Starting services...")

	// Start services
	if err := docker.ComposeUp(composePath, args...); err != nil {
		return fmt.Errorf("failed to start services: %w", err)
	}

	fmt.Println("✅ Services started successfully!")

	// Show service status
	if len(args) == 0 {
		fmt.Printf("   Started: %v\n", cfg.Services)
	} else {
		fmt.Printf("   Started: %v\n", args)
	}

	// Show connection info
	fmt.Println("\n📝 Connection Information:")
	servicesToShow := args
	if len(servicesToShow) == 0 {
		servicesToShow = cfg.Services
	}

	showConnectionInfo(cfg, servicesToShow)

	fmt.Println("\nManage services:")
	fmt.Println("  dockenv down     # Stop all services")
	fmt.Println("  dockenv restart  # Restart services")
	fmt.Println("  dockenv status   # Check service status")
	fmt.Println("  dockenv logs     # View service logs")

	return nil
}

func showConnectionInfo(cfg *config.Config, serviceNames []string) {
	for _, serviceName := range serviceNames {
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
}
