package config

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/yunier-rojas/minipaas/minipaas-cli/internal/core/entities"
	"github.com/yunier-rojas/minipaas/minipaas-cli/internal/core/usecases"
	"github.com/yunier-rojas/minipaas/minipaas-cli/internal/frameworks/logger"
)

type EnvironmentService struct{}

func NewEnvironmentService() *EnvironmentService {
	return &EnvironmentService{}
}

func (s *EnvironmentService) Apply(_ string, cfg entities.MinipaasConfig, verbose bool) {
	keys := make([]string, 0, len(cfg.Vars))
	for key := range cfg.Vars {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	for _, key := range keys {
		_ = os.Setenv(key, cfg.Vars[key])
	}

	_ = os.Unsetenv("DOCKER_CERT_PATH")
	_ = os.Unsetenv("DOCKER_TLS_VERIFY")

	if cfg.Docker.Host == "" {
		if verbose {
			logger.Info(s.describeVars(cfg.Vars, keys))
		}
		return
	}

	_ = os.Setenv("DOCKER_HOST", cfg.Docker.Host)

	if verbose {
		logger.Info(fmt.Sprintf(
			"%s\n   DOCKER_HOST=%s",
			s.describeVars(cfg.Vars, keys), cfg.Docker.Host))
	}
}

func (s *EnvironmentService) describeVars(vars map[string]string, keys []string) string {
	lines := []string{"Environment:"}
	for _, key := range keys {
		lines = append(lines, fmt.Sprintf("   %s=%s", key, vars[key]))
	}
	return strings.Join(lines, "\n")
}

var _ usecases.EnvironmentPort = (*EnvironmentService)(nil)
