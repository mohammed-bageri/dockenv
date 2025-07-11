package templates

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/mohammed-bageri/dockenv/internal/config"
	"github.com/mohammed-bageri/dockenv/internal/services"
)

type TemplateData struct {
	Port     int
	DataPath string
	Env      map[string]string
}

type ComposeData struct {
	Version  string
	Services map[string]TemplateData
	Networks bool
	Volumes  map[string]string
}

func GenerateDockerCompose(cfg *config.Config) error {
	composeFile := config.GetComposePath()

	file, err := os.Create(composeFile)
	if err != nil {
		return fmt.Errorf("failed to create compose file: %w", err)
	}
	defer file.Close()

	// Write header
	fmt.Fprintln(file, "version: '3.8'")
	fmt.Fprintln(file, "")
	fmt.Fprintln(file, "services:")

	// Generate services
	for _, serviceName := range cfg.Services {
		service, exists := services.GetService(serviceName)
		if !exists {
			return fmt.Errorf("unknown service: %s", serviceName)
		}

		port := cfg.Ports[serviceName]
		if port == 0 {
			port = service.DefaultPort
		}

		templateData := TemplateData{
			Port:     port,
			DataPath: cfg.DataPath,
			Env:      cfg.Env,
		}

		if err := generateServiceTemplate(file, service.Template, templateData); err != nil {
			return fmt.Errorf("failed to generate template for %s: %w", serviceName, err)
		}

		fmt.Fprintln(file, "")
	}

	// Add volumes section
	fmt.Fprintln(file, "volumes:")
	for _, serviceName := range cfg.Services {
		service, _ := services.GetService(serviceName)
		for _, volume := range service.Volumes {
			fmt.Fprintf(file, "  %s:\n", volume)
		}
	}

	return nil
}

func generateServiceTemplate(file *os.File, templateName string, data TemplateData) error {
	templatePath := filepath.Join("templates", templateName)

	content, err := os.ReadFile(templatePath)
	if err != nil {
		return fmt.Errorf("failed to read template %s: %w", templateName, err)
	}

	tmpl, err := template.New(templateName).Parse(string(content))
	if err != nil {
		return fmt.Errorf("failed to parse template %s: %w", templateName, err)
	}

	if err := tmpl.Execute(file, data); err != nil {
		return fmt.Errorf("failed to execute template %s: %w", templateName, err)
	}

	return nil
}

func GetEmbeddedTemplate(serviceName string) (string, error) {
	// Return embedded templates as fallback when template files don't exist
	templates := map[string]string{
		"mysql": `  mysql:
    image: mysql:8.0
    container_name: dockenv-mysql
    restart: unless-stopped
    environment:
      MYSQL_ROOT_PASSWORD: root
      MYSQL_DATABASE: dockenv
      MYSQL_USER: dockenv
      MYSQL_PASSWORD: password
    ports:
      - "{{.Port}}:3306"
    volumes:
      - {{.DataPath}}/mysql:/var/lib/mysql
    healthcheck:
      test: ["CMD", "mysqladmin", "ping", "-h", "localhost"]
      timeout: 20s
      retries: 10`,

		"postgres": `  postgres:
    image: postgres:15
    container_name: dockenv-postgres
    restart: unless-stopped
    environment:
      POSTGRES_DB: dockenv
      POSTGRES_USER: dockenv
      POSTGRES_PASSWORD: password
    ports:
      - "{{.Port}}:5432"
    volumes:
      - {{.DataPath}}/postgres:/var/lib/postgresql/data
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U dockenv"]
      interval: 30s
      timeout: 10s
      retries: 5`,

		"redis": `  redis:
    image: redis:7-alpine
    container_name: dockenv-redis
    restart: unless-stopped
    ports:
      - "{{.Port}}:6379"
    volumes:
      - {{.DataPath}}/redis:/data
    healthcheck:
      test: ["CMD", "redis-cli", "ping"]
      interval: 30s
      timeout: 10s
      retries: 5`,

		"mongodb": `  mongodb:
    image: mongo:7
    container_name: dockenv-mongodb
    restart: unless-stopped
    environment:
      MONGO_INITDB_ROOT_USERNAME: dockenv
      MONGO_INITDB_ROOT_PASSWORD: password
      MONGO_INITDB_DATABASE: dockenv
    ports:
      - "{{.Port}}:27017"
    volumes:
      - {{.DataPath}}/mongodb:/data/db
    healthcheck:
      test: ["CMD", "mongo", "--eval", "db.adminCommand('ping')"]
      interval: 30s
      timeout: 10s
      retries: 5`,

		"kafka": `  zookeeper:
    image: confluentinc/cp-zookeeper:latest
    container_name: dockenv-zookeeper
    restart: unless-stopped
    environment:
      ZOOKEEPER_CLIENT_PORT: 2181
      ZOOKEEPER_TICK_TIME: 2000
    volumes:
      - {{.DataPath}}/zookeeper:/var/lib/zookeeper/data

  kafka:
    image: confluentinc/cp-kafka:latest
    container_name: dockenv-kafka
    restart: unless-stopped
    depends_on:
      - zookeeper
    environment:
      KAFKA_BROKER_ID: 1
      KAFKA_ZOOKEEPER_CONNECT: zookeeper:2181
      KAFKA_ADVERTISED_LISTENERS: PLAINTEXT://localhost:{{.Port}}
      KAFKA_OFFSETS_TOPIC_REPLICATION_FACTOR: 1
    ports:
      - "{{.Port}}:9092"
    volumes:
      - {{.DataPath}}/kafka:/var/lib/kafka/data`,

		"elasticsearch": `  elasticsearch:
    image: docker.elastic.co/elasticsearch/elasticsearch:8.11.0
    container_name: dockenv-elasticsearch
    restart: unless-stopped
    environment:
      - discovery.type=single-node
      - xpack.security.enabled=false
      - "ES_JAVA_OPTS=-Xms512m -Xmx512m"
    ports:
      - "{{.Port}}:9200"
    volumes:
      - {{.DataPath}}/elasticsearch:/usr/share/elasticsearch/data
    healthcheck:
      test: ["CMD-SHELL", "curl -f http://localhost:9200/_cluster/health || exit 1"]
      interval: 30s
      timeout: 10s
      retries: 5`,

		"rabbitmq": `  rabbitmq:
    image: rabbitmq:3-management
    container_name: dockenv-rabbitmq
    restart: unless-stopped
    environment:
      RABBITMQ_DEFAULT_USER: dockenv
      RABBITMQ_DEFAULT_PASS: password
    ports:
      - "{{.Port}}:5672"
      - "15672:15672"
    volumes:
      - {{.DataPath}}/rabbitmq:/var/lib/rabbitmq
    healthcheck:
      test: ["CMD", "rabbitmq-diagnostics", "ping"]
      interval: 30s
      timeout: 10s
      retries: 5`,
	}

	templateStr, exists := templates[serviceName]
	if !exists {
		return "", fmt.Errorf("no embedded template found for service: %s", serviceName)
	}

	return templateStr, nil
}

func GenerateDockerComposeEmbedded(cfg *config.Config) error {
	composeFile := config.GetComposePath()

	file, err := os.Create(composeFile)
	if err != nil {
		return fmt.Errorf("failed to create compose file: %w", err)
	}
	defer file.Close()

	// Write header
	fmt.Fprintln(file, "version: '3.8'")
	fmt.Fprintln(file, "")
	fmt.Fprintln(file, "services:")

	// Generate services using embedded templates
	for _, serviceName := range cfg.Services {
		service, exists := services.GetService(serviceName)
		if !exists {
			return fmt.Errorf("unknown service: %s", serviceName)
		}

		port := cfg.Ports[serviceName]
		if port == 0 {
			port = service.DefaultPort
		}

		templateData := TemplateData{
			Port:     port,
			DataPath: cfg.DataPath,
			Env:      cfg.Env,
		}

		templateStr, err := GetEmbeddedTemplate(serviceName)
		if err != nil {
			return fmt.Errorf("failed to get embedded template for %s: %w", serviceName, err)
		}

		tmpl, err := template.New(serviceName).Parse(templateStr)
		if err != nil {
			return fmt.Errorf("failed to parse embedded template for %s: %w", serviceName, err)
		}

		if err := tmpl.Execute(file, templateData); err != nil {
			return fmt.Errorf("failed to execute embedded template for %s: %w", serviceName, err)
		}

		fmt.Fprintln(file, "")
	}

	// Add volumes section
	fmt.Fprintln(file, "volumes:")
	volumes := make(map[string]bool)
	for _, serviceName := range cfg.Services {
		service, _ := services.GetService(serviceName)
		for _, volume := range service.Volumes {
			if !volumes[volume] {
				fmt.Fprintf(file, "  %s:\n", volume)
				volumes[volume] = true
			}
		}
	}

	return nil
}

func UpdateDockerCompose(cfg *config.Config, servicesToAdd []string, servicesToRemove []string) error {
	// Add new services
	for _, serviceName := range servicesToAdd {
		if !contains(cfg.Services, serviceName) {
			cfg.Services = append(cfg.Services, serviceName)

			// Add default port if not set
			if cfg.Ports[serviceName] == 0 {
				if service, exists := services.GetService(serviceName); exists {
					cfg.Ports[serviceName] = service.DefaultPort
				}
			}

			// Add env vars
			if service, exists := services.GetService(serviceName); exists {
				for key, value := range service.EnvVars {
					if cfg.Env[key] == "" {
						cfg.Env[key] = value
					}
				}
			}
		}
	}

	// Remove services
	for _, serviceName := range servicesToRemove {
		cfg.Services = removeString(cfg.Services, serviceName)
		delete(cfg.Ports, serviceName)

		// Remove env vars (only those specific to this service)
		if service, exists := services.GetService(serviceName); exists {
			for key := range service.EnvVars {
				delete(cfg.Env, key)
			}
		}
	}

	// Regenerate Docker Compose file
	return GenerateDockerComposeEmbedded(cfg)
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func removeString(slice []string, item string) []string {
	var result []string
	for _, s := range slice {
		if s != item {
			result = append(result, s)
		}
	}
	return result
}
