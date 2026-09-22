package config

import (
	"os"
	"testing"

	"github.com/yunier-rojas/minipaas/minipaas-cli/internal/core/entities"
)

func unsetDockerEnv() {
	_ = os.Unsetenv("DOCKER_CERT_PATH")
	_ = os.Unsetenv("DOCKER_HOST")
	_ = os.Unsetenv("DOCKER_TLS_VERIFY")
}

func TestEnvironmentService_ExportsVarsVerbatim(t *testing.T) {
	_ = os.Unsetenv("MINIPAAS_DEPLOY_VERSION")
	_ = os.Unsetenv("LOG_LEVEL")

	cfg := entities.MinipaasConfig{Vars: map[string]string{
		"MINIPAAS_DEPLOY_VERSION": "1.2.3",
		"LOG_LEVEL":               "debug",
	}}
	NewEnvironmentService().Apply("/tmp/env", cfg, false)

	if os.Getenv("MINIPAAS_DEPLOY_VERSION") != "1.2.3" || os.Getenv("LOG_LEVEL") != "debug" {
		t.Fatalf("vars not exported: version=%q log=%q", os.Getenv("MINIPAAS_DEPLOY_VERSION"), os.Getenv("LOG_LEVEL"))
	}
}

func TestEnvironmentService_EmptyHostIsLocal(t *testing.T) {
	unsetDockerEnv()

	cfg := entities.MinipaasConfig{Vars: map[string]string{"MINIPAAS_DEPLOY_VERSION": "1.2.3"}}
	NewEnvironmentService().Apply("/tmp/env", cfg, true)

	if os.Getenv("DOCKER_TLS_VERIFY") != "" || os.Getenv("DOCKER_HOST") != "" || os.Getenv("DOCKER_CERT_PATH") != "" {
		t.Fatalf("docker envs should be empty when host is omitted")
	}
}

func TestEnvironmentService_RemoteHostSetsDockerHost(t *testing.T) {
	unsetDockerEnv()
	_ = os.Setenv("DOCKER_CERT_PATH", "/stale")
	_ = os.Setenv("DOCKER_TLS_VERIFY", "1")

	cfg := entities.MinipaasConfig{
		Docker: entities.DockerConfig{Host: "ssh://deploy@example.com"},
		Vars:   map[string]string{"MINIPAAS_DEPLOY_VERSION": "9.9.9"},
	}
	NewEnvironmentService().Apply("/tmp/env", cfg, true)

	if os.Getenv("DOCKER_HOST") != "ssh://deploy@example.com" {
		t.Fatalf("host not set")
	}
	if os.Getenv("DOCKER_CERT_PATH") != "" || os.Getenv("DOCKER_TLS_VERIFY") != "" {
		t.Fatalf("tls envs should be cleared")
	}
	if os.Getenv("MINIPAAS_DEPLOY_VERSION") != "9.9.9" {
		t.Fatalf("vars not exported alongside docker env")
	}
}
