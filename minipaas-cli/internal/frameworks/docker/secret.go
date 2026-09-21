package docker

import (
	"github.com/yunier-rojas/minipaas/minipaas-cli/internal/core/usecases"
	"github.com/yunier-rojas/minipaas/minipaas-cli/internal/frameworks/process"
)

type SecretService struct {
	runner process.CommandRunner
}

func NewSecretService(runner process.CommandRunner) *SecretService {
	return &SecretService{runner: runner}
}

func (s *SecretService) Exists(name string) bool {
	return s.runner.Run([]string{"docker", "secret", "inspect", name}, false) == nil
}

func (s *SecretService) Create(name string, content []byte, verbose bool) error {
	return s.runner.RunWithInput([]string{"docker", "secret", "create", name, "-"}, content, verbose)
}

var _ usecases.DockerSecretPort = (*SecretService)(nil)
