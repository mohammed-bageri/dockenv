package docker

import (
	"os"
	"os/exec"
	"strings"
)

type DockerInfo struct {
	DockerInstalled        bool
	DockerComposeInstalled bool
	DockerVersion          string
	ComposeVersion         string
	DockerRunning          bool
}

func CheckDocker() (*DockerInfo, error) {
	info := &DockerInfo{}

	// Check if Docker is installed
	if dockerCmd, err := exec.LookPath("docker"); err == nil && dockerCmd != "" {
		info.DockerInstalled = true

		// Get Docker version
		if out, err := exec.Command("docker", "--version").Output(); err == nil {
			info.DockerVersion = strings.TrimSpace(string(out))
		}

		// Check if Docker daemon is running
		if err := exec.Command("docker", "info").Run(); err == nil {
			info.DockerRunning = true
		}
	}

	// Check if Docker Compose is installed
	if composeCmd, err := exec.LookPath("docker-compose"); err == nil && composeCmd != "" {
		info.DockerComposeInstalled = true

		// Get Docker Compose version
		if out, err := exec.Command("docker-compose", "--version").Output(); err == nil {
			info.ComposeVersion = strings.TrimSpace(string(out))
		}
	} else {
		// Check for docker compose (newer syntax)
		if err := exec.Command("docker", "compose", "version").Run(); err == nil {
			info.DockerComposeInstalled = true
			if out, err := exec.Command("docker", "compose", "version").Output(); err == nil {
				info.ComposeVersion = strings.TrimSpace(string(out))
			}
		}
	}

	return info, nil
}

func (d *DockerInfo) IsReady() bool {
	return d.DockerInstalled && d.DockerComposeInstalled && d.DockerRunning
}

func (d *DockerInfo) GetInstallInstructions() string {
	var instructions []string

	if !d.DockerInstalled {
		instructions = append(instructions, "Docker is not installed. Please install Docker from https://docs.docker.com/get-docker/")
	} else if !d.DockerRunning {
		instructions = append(instructions, "Docker daemon is not running. Please start Docker.")
	}

	if !d.DockerComposeInstalled {
		instructions = append(instructions, "Docker Compose is not installed. Please install Docker Compose from https://docs.docker.com/compose/install/")
	}

	if len(instructions) == 0 {
		return "Docker and Docker Compose are properly installed and running!"
	}

	return strings.Join(instructions, "\n")
}

func RunCompose(args ...string) error {
	// Try docker compose first (newer syntax)
	cmd := exec.Command("docker", append([]string{"compose"}, args...)...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		// Fallback to docker-compose
		cmd = exec.Command("docker-compose", args...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		return cmd.Run()
	}

	return nil
}

func ComposeUp(file string, services ...string) error {
	args := []string{"-f", file, "up", "-d"}
	if len(services) > 0 {
		args = append(args, services...)
	}
	return RunCompose(args...)
}

func ComposeDown(file string) error {
	return RunCompose("-f", file, "down")
}

func ComposeRestart(file string, services ...string) error {
	args := []string{"-f", file, "restart"}
	if len(services) > 0 {
		args = append(args, services...)
	}
	return RunCompose(args...)
}

func ComposeStatus(file string) error {
	return RunCompose("-f", file, "ps")
}

func ComposeLogs(file string, services ...string) error {
	args := []string{"-f", file, "logs", "-f"}
	if len(services) > 0 {
		args = append(args, services...)
	}
	return RunCompose(args...)
}

func ComposeValidate(file string) error {
	return RunCompose("-f", file, "config", "--quiet")
}
