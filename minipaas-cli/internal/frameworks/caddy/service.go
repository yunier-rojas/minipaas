package caddy

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/cockroachdb/errors"

	"github.com/yunier-rojas/minipaas/minipaas-cli/internal/core/usecases"
	"github.com/yunier-rojas/minipaas/minipaas-cli/internal/frameworks/logger"
	"github.com/yunier-rojas/minipaas/minipaas-cli/internal/frameworks/process"
)

const configFile = "caddy.json"

type Service struct {
	runner process.CommandRunner
}

func NewService(runner process.CommandRunner) *Service {
	return &Service{runner: runner}
}

func (s *Service) Apply(env string, namespace string, verbose bool) error {
	file := filepath.Join(env, configFile)
	payload, err := os.ReadFile(file)
	if err != nil {
		return errors.Wrap(err, "read caddy config")
	}

	containerID, err := s.containerID(namespace)
	if err != nil {
		return errors.Wrap(err, "find caddy container")
	}

	args := []string{
		"docker", "exec", "-i", containerID,
		"/usr/bin/wget", "-O", "-", "-q",
		"--header=Content-Type: application/json",
		"--post-data=" + string(payload),
		"http://127.0.0.1:2019/load",
	}
	if err = s.runner.Run(args, verbose); err != nil {
		return errors.Wrap(err, "update caddy routing")
	}

	logger.Success("Routing updated")
	return nil
}

func (s *Service) containerID(namespace string) (string, error) {
	serviceName := namespace + "_caddy"

	output, err := s.runner.RunOutput([]string{
		"docker", "ps",
		"--filter", "label=com.docker.swarm.service.name=" + serviceName,
		"--format", "{{.ID}}",
	}, false)
	if err != nil {
		return "", errors.Wrap(err, "docker ps")
	}

	ids := strings.Fields(output)
	if len(ids) == 0 {
		return "", errors.New("no running container found for service " + serviceName)
	}

	return ids[0], nil
}

var _ usecases.RoutingPort = (*Service)(nil)
