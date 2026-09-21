package docker

import (
	"github.com/yunier-rojas/minipaas/minipaas-cli/internal/core/usecases"
	"github.com/yunier-rojas/minipaas/minipaas-cli/internal/frameworks/process"
)

type ConfigService struct {
	runner process.CommandRunner
}

func NewConfigService(runner process.CommandRunner) *ConfigService {
	return &ConfigService{runner: runner}
}

func (s *ConfigService) Exists(name string) bool {
	return s.runner.Run([]string{"docker", "config", "inspect", name}, false) == nil
}

func (s *ConfigService) Create(name string, content []byte, verbose bool) error {
	return s.runner.RunWithInput([]string{"docker", "config", "create", name, "-"}, content, verbose)
}

var _ usecases.DockerConfigPort = (*ConfigService)(nil)
