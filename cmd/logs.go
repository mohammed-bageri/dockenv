package cmd

import (
	"fmt"

	"github.com/mohammed-bageri/dockenv/internal/config"
	"github.com/mohammed-bageri/dockenv/internal/docker"
	"github.com/mohammed-bageri/dockenv/internal/utils"

	"github.com/spf13/cobra"
)

var logsCmd = &cobra.Command{
	Use:   "logs [service...]",
	Short: "Show service logs",
	Long: `Show logs from development services.
You can optionally specify which services to show logs for,
otherwise logs from all services will be displayed.

Examples:
  dockenv logs           # Show logs from all services
  dockenv logs mysql     # Show only MySQL logs
  dockenv logs -f mysql  # Follow MySQL logs`,
	RunE: runLogs,
}

var (
	followLogsFlag bool
	tailFlag       string
)

func init() {
	rootCmd.AddCommand(logsCmd)

	logsCmd.Flags().BoolVarP(&followLogsFlag, "follow", "f", false, "Follow log output")
	logsCmd.Flags().StringVarP(&tailFlag, "tail", "t", "100", "Number of lines to show from end of logs")
}

func runLogs(cmd *cobra.Command, args []string) error {
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
	if len(args) > 0 {
		cfg, err := utils.LoadConfig()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		// Validate requested services
		for _, serviceName := range args {
			if !utils.Contains(cfg.Services, serviceName) {
				return fmt.Errorf("service '%s' not configured. Available services: %v", serviceName, cfg.Services)
			}
		}
	}

	fmt.Println("📋 Service Logs:")
	if len(args) == 0 {
		fmt.Println("   Showing logs from all services")
	} else {
		fmt.Printf("   Showing logs from: %v\n", args)
	}

	if followLogsFlag {
		fmt.Println("   Press Ctrl+C to stop following logs")
	}
	fmt.Println()

	// Show logs
	if followLogsFlag {
		return docker.ComposeLogs(composePath, args...)
	} else {
		// For non-follow mode, use docker-compose logs with tail
		logArgs := []string{"-f", composePath, "logs", "--tail", tailFlag}
		if len(args) > 0 {
			logArgs = append(logArgs, args...)
		}
		return docker.RunCompose(logArgs...)
	}
}
