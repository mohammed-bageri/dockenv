package services

import (
	"fmt"
	"strings"
)

type Service struct {
	Name        string
	DisplayName string
	Description string
	DefaultPort int
	Template    string
	Volumes     []string
	EnvVars     map[string]string
}

var AvailableServices = map[string]Service{
	"mysql": {
		Name:        "mysql",
		DisplayName: "MySQL",
		Description: "MySQL Database Server",
		DefaultPort: 3306,
		Template:    "mysql.yaml",
		Volumes:     []string{"mysql_data"},
		EnvVars: map[string]string{
			"DB_CONNECTION": "mysql",
			"DB_HOST":       "127.0.0.1",
			"DB_PORT":       "3306",
			"DB_DATABASE":   "dockenv",
			"DB_USERNAME":   "dockenv",
			"DB_PASSWORD":   "password",
		},
	},
	"postgres": {
		Name:        "postgres",
		DisplayName: "PostgreSQL",
		Description: "PostgreSQL Database Server",
		DefaultPort: 5432,
		Template:    "postgres.yaml",
		Volumes:     []string{"postgres_data"},
		EnvVars: map[string]string{
			"DB_CONNECTION": "pgsql",
			"DB_HOST":       "127.0.0.1",
			"DB_PORT":       "5432",
			"DB_DATABASE":   "dockenv",
			"DB_USERNAME":   "dockenv",
			"DB_PASSWORD":   "password",
		},
	},
	"redis": {
		Name:        "redis",
		DisplayName: "Redis",
		Description: "Redis In-Memory Data Store",
		DefaultPort: 6379,
		Template:    "redis.yaml",
		Volumes:     []string{"redis_data"},
		EnvVars: map[string]string{
			"REDIS_HOST":     "127.0.0.1",
			"REDIS_PORT":     "6379",
			"REDIS_PASSWORD": "",
		},
	},
	"mongodb": {
		Name:        "mongodb",
		DisplayName: "MongoDB",
		Description: "MongoDB NoSQL Database",
		DefaultPort: 27017,
		Template:    "mongodb.yaml",
		Volumes:     []string{"mongodb_data"},
		EnvVars: map[string]string{
			"MONGO_HOST":     "127.0.0.1",
			"MONGO_PORT":     "27017",
			"MONGO_DATABASE": "dockenv",
			"MONGO_USERNAME": "dockenv",
			"MONGO_PASSWORD": "password",
		},
	},
	"kafka": {
		Name:        "kafka",
		DisplayName: "Apache Kafka",
		Description: "Apache Kafka Message Broker",
		DefaultPort: 9092,
		Template:    "kafka.yaml",
		Volumes:     []string{"kafka_data", "zookeeper_data"},
		EnvVars: map[string]string{
			"KAFKA_HOST": "127.0.0.1",
			"KAFKA_PORT": "9092",
		},
	},
	"elasticsearch": {
		Name:        "elasticsearch",
		DisplayName: "Elasticsearch",
		Description: "Elasticsearch Search Engine",
		DefaultPort: 9200,
		Template:    "elasticsearch.yaml",
		Volumes:     []string{"elasticsearch_data"},
		EnvVars: map[string]string{
			"ELASTICSEARCH_HOST": "127.0.0.1",
			"ELASTICSEARCH_PORT": "9200",
		},
	},
	"rabbitmq": {
		Name:        "rabbitmq",
		DisplayName: "RabbitMQ",
		Description: "RabbitMQ Message Broker",
		DefaultPort: 5672,
		Template:    "rabbitmq.yaml",
		Volumes:     []string{"rabbitmq_data"},
		EnvVars: map[string]string{
			"RABBITMQ_HOST":     "127.0.0.1",
			"RABBITMQ_PORT":     "5672",
			"RABBITMQ_USERNAME": "dockenv",
			"RABBITMQ_PASSWORD": "password",
		},
	},
}

// Profiles define common service bundles
var Profiles = map[string][]string{
	"laravel": {"mysql", "redis"},
	"node":    {"postgres", "redis"},
	"django":  {"postgres", "redis"},
	"rails":   {"postgres", "redis"},
	"spring":  {"mysql", "kafka"},
	"full":    {"mysql", "postgres", "redis", "mongodb", "kafka"},
}

func GetService(name string) (Service, bool) {
	service, exists := AvailableServices[name]
	return service, exists
}

func GetServiceNames() []string {
	names := make([]string, 0, len(AvailableServices))
	for name := range AvailableServices {
		names = append(names, name)
	}
	return names
}

func GetServiceDisplayNames() []string {
	names := make([]string, 0, len(AvailableServices))
	for _, service := range AvailableServices {
		names = append(names, fmt.Sprintf("%s - %s", service.DisplayName, service.Description))
	}
	return names
}

func ValidateServices(services []string) error {
	for _, serviceName := range services {
		if _, exists := AvailableServices[serviceName]; !exists {
			return fmt.Errorf("unknown service: %s. Available services: %s",
				serviceName, strings.Join(GetServiceNames(), ", "))
		}
	}
	return nil
}

func GetProfileServices(profileName string) ([]string, bool) {
	services, exists := Profiles[profileName]
	return services, exists
}

func GetProfileNames() []string {
	names := make([]string, 0, len(Profiles))
	for name := range Profiles {
		names = append(names, name)
	}
	return names
}
