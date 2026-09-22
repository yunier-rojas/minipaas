package deploy

import (
	"context"
	"fmt"
	"path/filepath"
	"strconv"

	"github.com/cockroachdb/errors"
	"github.com/compose-spec/compose-go/v2/cli"
	"github.com/compose-spec/compose-go/v2/types"

	"github.com/yunier-rojas/minipaas/minipaas-cli/internal/core/usecases"
	"github.com/yunier-rojas/minipaas/minipaas-cli/internal/frameworks/logger"
	"github.com/yunier-rojas/minipaas/minipaas-cli/internal/frameworks/process"
)

type Service struct {
	runner process.CommandRunner
}

func NewService(runner process.CommandRunner) *Service {
	return &Service{runner: runner}
}

func (s *Service) Build(files []string, registry string, namespace string, verbose bool) error {
	project, err := loadDeployProject(files, namespace)
	if err != nil {
		return errors.Wrap(err, "load deploy project")
	}

	var errs []error
	for name, svc := range project.Services {
		if svc.Build == nil {
			continue
		}
		if err = s.runner.Run(buildCommandFromService(svc), verbose); err != nil {
			logger.Error(fmt.Sprintf("Failed to build %s", name), err)
			errs = append(errs, errors.Wrapf(err, "build %s", name))
			continue
		}
		logger.Success(fmt.Sprintf("Built %s: %s", name, svc.Image))

		if registry == "" || svc.Image == "" {
			continue
		}
		if err = s.runner.Run([]string{"docker", "push", svc.Image}, verbose); err != nil {
			logger.Error(fmt.Sprintf("Failed to push %s", name), err)
			errs = append(errs, errors.Wrapf(err, "push %s", name))
			continue
		}
		logger.Success(fmt.Sprintf("Pushed %s: %s", name, svc.Image))
	}

	return errors.Join(errs...)
}

func (s *Service) Rollout(files []string, namespace string, verbose bool) error {
	args := []string{"docker", "stack", "deploy", "--with-registry-auth"}
	for _, file := range files {
		args = append(args, "-c", file)
	}
	args = append(args, namespace)

	if err := s.runner.Run(args, verbose); err != nil {
		return errors.Wrap(err, "deploy stack")
	}

	logger.Success("Deployment successful")
	return nil
}

func (s *Service) Canary(files []string, services []string, replicas int, namespace string, verbose bool) error {
	project, err := loadDeployProject(files, namespace)
	if err != nil {
		return errors.Wrap(err, "load deploy project")
	}

	var errs []error
	for _, name := range services {
		svc, err := project.GetService(name)
		if err != nil {
			logger.Error(fmt.Sprintf("Failed to get service %s", name), err)
			errs = append(errs, errors.Wrapf(err, "get service %s", name))
			continue
		}

		args := []string{
			"docker", "service", "update",
			"--with-registry-auth",
			"--replicas", strconv.Itoa(replicas),
			"--image", svc.Image,
			namespace + "_" + name,
		}
		if err = s.runner.Run(args, verbose); err != nil {
			logger.Error(fmt.Sprintf("Failed to do canary release for service %s", name), err)
			errs = append(errs, errors.Wrapf(err, "canary %s", name))
			continue
		}
		logger.Success(fmt.Sprintf("Canary released %s: %s", name, svc.Image))
	}

	return errors.Join(errs...)
}

func loadDeployProject(files []string, namespace string) (*types.Project, error) {
	options, err := cli.NewProjectOptions(
		files,
		cli.WithName(namespace),
		cli.WithResolvedPaths(false),
		cli.WithConsistency(false),
		cli.WithOsEnv,
	)
	if err != nil {
		return nil, err
	}

	return options.LoadProject(context.Background())
}

func buildCommandFromService(svc types.ServiceConfig) []string {
	cmd := []string{"docker", "build", "--network", "host"}

	if svc.Image != "" {
		cmd = append(cmd, "-t", svc.Image)
	}

	if svc.Build == nil {
		return append(cmd, ".")
	}

	context := svc.Build.Context
	if context == "" {
		context = "."
	}

	if svc.Build.Dockerfile != "" {
		cmd = append(cmd, "-f", filepath.Join(context, svc.Build.Dockerfile))
	}

	for key, valPtr := range svc.Build.Args {
		if valPtr != nil {
			cmd = append(cmd, "--build-arg", fmt.Sprintf("%s=%s", key, *valPtr))
		} else {
			cmd = append(cmd, "--build-arg", fmt.Sprintf("%s=", key))
		}
	}

	return append(cmd, context)
}

var _ usecases.DeployPort = (*Service)(nil)
