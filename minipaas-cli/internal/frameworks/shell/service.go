package shell

import (
	"os"
	"os/exec"

	"github.com/yunier-rojas/minipaas/minipaas-cli/internal/core/usecases"
	"github.com/yunier-rojas/minipaas/minipaas-cli/internal/frameworks/logger"
)

type Service struct{}

func NewService() *Service {
	return &Service{}
}

func (s *Service) Launch() error {
	sh := os.Getenv("SHELL")
	if sh == "" {
		sh = "/bin/sh"
	}
	logger.Info("Launching shell: " + sh)

	cmd := exec.Command(sh)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

var _ usecases.ShellPort = (*Service)(nil)
