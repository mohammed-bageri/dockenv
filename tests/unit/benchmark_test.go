package unit

import (
	"testing"

	"github.com/mohammed-bageri/dockenv/internal/config"
	"github.com/mohammed-bageri/dockenv/internal/services"
	"github.com/mohammed-bageri/dockenv/internal/templates"
	"github.com/mohammed-bageri/dockenv/internal/utils"
)

func BenchmarkLoadConfig(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = utils.LoadConfig()
	}
}

func BenchmarkGetServiceNames(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = services.GetServiceNames()
	}
}

func BenchmarkValidateServices(b *testing.B) {
	testServices := []string{"mysql", "redis", "postgres"}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = services.ValidateServices(testServices)
	}
}

func BenchmarkGenerateDockerComposeEmbedded(b *testing.B) {
	cfg := &config.Config{
		Version:  "1.0",
		Services: []string{"mysql", "redis"},
		Ports: map[string]int{
			"mysql": 3306,
			"redis": 6379,
		},
		Env:      map[string]string{},
		Volumes:  map[string]string{},
		DataPath: "/tmp/test",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = templates.GenerateDockerComposeEmbedded(cfg)
	}
}

func BenchmarkGetEmbeddedTemplate(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = templates.GetEmbeddedTemplate("mysql")
	}
}
