package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "dockenv",
	Short: "Manage local development environments with Docker Compose",
	Long: `dockenv is a CLI tool that helps you set up local development environments
using Docker Compose. It provides templates for common services like MySQL,
PostgreSQL, Redis, MongoDB, and Kafka without requiring you to install them
directly on your system.

Example usage:
  dockenv init         # Interactive setup
  dockenv up           # Start services
  dockenv down         # Stop services
  dockenv add mysql    # Add a service
  dockenv remove redis # Remove a service`,
	Version: "0.2.0",
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true
}
