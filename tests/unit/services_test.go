package unit

import (
	"testing"

	"github.com/mohammed-bageri/dockenv/internal/services"
)

func TestGetServiceNames(t *testing.T) {
	names := services.GetServiceNames()

	expectedServices := []string{"mysql", "postgres", "redis", "mongodb", "kafka", "elasticsearch", "rabbitmq"}

	if len(names) < len(expectedServices) {
		t.Errorf("Expected at least %d services, got %d", len(expectedServices), len(names))
	}

	for _, expected := range expectedServices {
		found := false
		for _, name := range names {
			if name == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected service %s not found in service names", expected)
		}
	}
}

func TestValidateServices(t *testing.T) {
	tests := []struct {
		name     string
		services []string
		hasError bool
	}{
		{"valid services", []string{"mysql", "redis"}, false},
		{"invalid service", []string{"mysql", "invalid"}, true},
		{"empty slice", []string{}, false},
		{"all invalid", []string{"invalid1", "invalid2"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := services.ValidateServices(tt.services)
			if tt.hasError && err == nil {
				t.Errorf("ValidateServices(%v) expected error, got nil", tt.services)
			}
			if !tt.hasError && err != nil {
				t.Errorf("ValidateServices(%v) unexpected error: %v", tt.services, err)
			}
		})
	}
}

func TestGetService(t *testing.T) {
	tests := []struct {
		name        string
		serviceName string
		expectFound bool
	}{
		{"get mysql", "mysql", true},
		{"get postgres", "postgres", true},
		{"get invalid", "invalid", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, found := services.GetService(tt.serviceName)
			if found != tt.expectFound {
				t.Errorf("GetService(%v) found = %v, want %v", tt.serviceName, found, tt.expectFound)
			}

			if found {
				if service.Name != tt.serviceName {
					t.Errorf("GetService(%v) name = %v, want %v", tt.serviceName, service.Name, tt.serviceName)
				}
				if service.DefaultPort == 0 {
					t.Errorf("GetService(%v) should have a default port", tt.serviceName)
				}
				if service.Template == "" {
					t.Errorf("GetService(%v) should have a template", tt.serviceName)
				}
			}
		})
	}
}

func TestServiceStructure(t *testing.T) {
	for serviceName, service := range services.AvailableServices {
		t.Run("service_"+serviceName, func(t *testing.T) {
			if service.Name != serviceName {
				t.Errorf("Service %s: name field mismatch, got %s", serviceName, service.Name)
			}
			if service.DisplayName == "" {
				t.Errorf("Service %s: DisplayName should not be empty", serviceName)
			}
			if service.Description == "" {
				t.Errorf("Service %s: Description should not be empty", serviceName)
			}
			if service.DefaultPort == 0 {
				t.Errorf("Service %s: DefaultPort should be set", serviceName)
			}
			if service.Template == "" {
				t.Errorf("Service %s: Template should not be empty", serviceName)
			}
		})
	}
}

func TestGetProfileServices(t *testing.T) {
	tests := []struct {
		name        string
		profileName string
		expectFound bool
		minServices int
	}{
		{"laravel profile", "laravel", true, 2},
		{"node profile", "node", true, 2},
		{"full profile", "full", true, 5},
		{"invalid profile", "invalid", false, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			serviceList, found := services.GetProfileServices(tt.profileName)
			if found != tt.expectFound {
				t.Errorf("GetProfileServices(%v) found = %v, want %v", tt.profileName, found, tt.expectFound)
			}

			if found && len(serviceList) < tt.minServices {
				t.Errorf("GetProfileServices(%v) returned %d services, expected at least %d",
					tt.profileName, len(serviceList), tt.minServices)
			}
		})
	}
}

func TestGetProfileNames(t *testing.T) {
	names := services.GetProfileNames()

	expectedProfiles := []string{"laravel", "node", "django", "rails", "spring", "full"}

	if len(names) < len(expectedProfiles) {
		t.Errorf("Expected at least %d profiles, got %d", len(expectedProfiles), len(names))
	}

	for _, expected := range expectedProfiles {
		found := false
		for _, name := range names {
			if name == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected profile %s not found in profile names", expected)
		}
	}
}
