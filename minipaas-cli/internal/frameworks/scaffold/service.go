package scaffold

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/yunier-rojas/minipaas/minipaas-cli/internal/core/entities"
	"github.com/yunier-rojas/minipaas/minipaas-cli/internal/core/usecases"
	"github.com/yunier-rojas/minipaas/minipaas-cli/internal/frameworks/logger"
)

//go:embed assets/*
var assets embed.FS

type Service struct{}

func NewService() *Service {
	return &Service{}
}

var baseAssets = []struct {
	target string
	source string
}{
	{"compose.common.yml", "common.yml"},
	{"compose.postgres.yml", "postgres.yaml"},
	{"compose.caddy.yml", "caddy.yml"},
	{"caddy.json", "caddy.json"},
}

func (s *Service) WriteBase(env string, namespace string) (entities.BaseFiles, error) {
	if err := os.MkdirAll(env, 0755); err != nil {
		return entities.BaseFiles{}, err
	}

	for _, asset := range baseAssets {
		raw, err := assets.ReadFile("assets/" + asset.source)
		if err != nil {
			return entities.BaseFiles{}, fmt.Errorf("embedded asset %s: %w", asset.source, err)
		}

		content := namespaceNetworks(string(raw), namespace)

		target := filepath.Join(env, asset.target)
		if err = os.WriteFile(target, []byte(content), 0644); err != nil {
			return entities.BaseFiles{}, err
		}

		logger.Info("Wrote file: " + target)
	}

	return entities.BaseFiles{
		Common:    "compose.common.yml",
		Postgres:  "compose.postgres.yml",
		Caddy:     "compose.caddy.yml",
		CaddyJSON: "caddy.json",
	}, nil
}

func namespaceNetworks(content, namespace string) string {
	content = strings.ReplaceAll(content, "minipaas_network", namespace+"_network")
	content = strings.ReplaceAll(content, "minipaas_public", namespace+"_public")
	return content
}

var _ usecases.ScaffoldPort = (*Service)(nil)
