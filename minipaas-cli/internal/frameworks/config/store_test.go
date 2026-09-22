package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yunier-rojas/minipaas/minipaas-cli/internal/core/entities"
)

func TestStoreService_Load(t *testing.T) {
	dir := t.TempDir()
	content := "compose:\n  files:\n    - compose.common.yml\n    - compose.apps.yaml\ndocker:\n  host: ssh://deploy@1.2.3.4\nvars:\n  MINIPAAS_DEPLOY_VERSION: 2.3.4\n"
	if err := os.WriteFile(filepath.Join(dir, entities.ConfigFile), []byte(content), 0644); err != nil {
		t.Fatalf("write config: %v", err)
	}

	cfg, err := NewStoreService().Load(dir)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(cfg.Compose.Files) != 2 || cfg.Compose.Files[0] != "compose.common.yml" {
		t.Fatalf("compose files mismatch: %#v", cfg.Compose.Files)
	}
	if cfg.Docker.Host != "ssh://deploy@1.2.3.4" {
		t.Fatalf("docker mismatch: %#v", cfg.Docker)
	}
	if cfg.Vars["MINIPAAS_DEPLOY_VERSION"] != "2.3.4" {
		t.Fatalf("vars mismatch: %#v", cfg.Vars)
	}
}

func TestStoreService_LoadMinimal(t *testing.T) {
	dir := t.TempDir()
	content := "compose:\n  files:\n    - compose.apps.yaml\n"
	if err := os.WriteFile(filepath.Join(dir, entities.ConfigFile), []byte(content), 0644); err != nil {
		t.Fatalf("write config: %v", err)
	}

	cfg, err := NewStoreService().Load(dir)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if cfg.Docker.Host != "" || len(cfg.Vars) != 0 {
		t.Fatalf("expected zero optional fields: %#v", cfg)
	}
}

func TestStoreService_LoadMissing(t *testing.T) {
	if _, err := NewStoreService().Load(t.TempDir()); err == nil {
		t.Fatalf("expected error for missing config")
	}
}
