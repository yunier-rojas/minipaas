package runtime

import (
	"go.uber.org/dig"

	"github.com/yunier-rojas/minipaas/minipaas-cli/internal/core/usecases"
	"github.com/yunier-rojas/minipaas/minipaas-cli/internal/frameworks/caddy"
	"github.com/yunier-rojas/minipaas/minipaas-cli/internal/frameworks/compose"
	"github.com/yunier-rojas/minipaas/minipaas-cli/internal/frameworks/config"
	"github.com/yunier-rojas/minipaas/minipaas-cli/internal/frameworks/deploy"
	"github.com/yunier-rojas/minipaas/minipaas-cli/internal/frameworks/docker"
	"github.com/yunier-rojas/minipaas/minipaas-cli/internal/frameworks/files"
	"github.com/yunier-rojas/minipaas/minipaas-cli/internal/frameworks/process"
	"github.com/yunier-rojas/minipaas/minipaas-cli/internal/frameworks/scaffold"
	"github.com/yunier-rojas/minipaas/minipaas-cli/internal/frameworks/shell"
)

type provider struct {
	constructor interface{}
	options     []dig.ProvideOption
}

func newContainer() (*dig.Container, error) {
	container := dig.New()

	providers := []provider{
		{constructor: process.NewRunner, options: []dig.ProvideOption{dig.As(new(process.CommandRunner))}},

		{constructor: config.NewStoreService, options: []dig.ProvideOption{dig.As(new(usecases.ConfigStorePort))}},
		{constructor: config.NewEnvironmentService, options: []dig.ProvideOption{dig.As(new(usecases.EnvironmentPort))}},
		{constructor: files.NewReaderService, options: []dig.ProvideOption{dig.As(new(usecases.InputPort))}},
		{constructor: compose.NewManagerService, options: []dig.ProvideOption{dig.As(new(usecases.ComposePort))}},
		{constructor: scaffold.NewService, options: []dig.ProvideOption{dig.As(new(usecases.ScaffoldPort))}},
		{constructor: docker.NewSecretService, options: []dig.ProvideOption{dig.As(new(usecases.DockerSecretPort))}},
		{constructor: docker.NewConfigService, options: []dig.ProvideOption{dig.As(new(usecases.DockerConfigPort))}},
		{constructor: shell.NewService, options: []dig.ProvideOption{dig.As(new(usecases.ShellPort))}},
		{constructor: deploy.NewService, options: []dig.ProvideOption{dig.As(new(usecases.DeployPort))}},
		{constructor: caddy.NewService, options: []dig.ProvideOption{dig.As(new(usecases.RoutingPort))}},
		{constructor: caddy.NewConfigService, options: []dig.ProvideOption{dig.As(new(usecases.RoutePort))}},

		{constructor: usecases.NewSecretCreateInteractor, options: []dig.ProvideOption{dig.As(new(usecases.SecretCreateCase))}},
		{constructor: usecases.NewConfigCreateInteractor, options: []dig.ProvideOption{dig.As(new(usecases.ConfigCreateCase))}},
		{constructor: usecases.NewShellInteractor, options: []dig.ProvideOption{dig.As(new(usecases.ShellCase))}},
		{constructor: usecases.NewDeployBuildInteractor, options: []dig.ProvideOption{dig.As(new(usecases.DeployBuildCase))}},
		{constructor: usecases.NewDeployRolloutInteractor, options: []dig.ProvideOption{dig.As(new(usecases.DeployRolloutCase))}},
		{constructor: usecases.NewDeployCanaryInteractor, options: []dig.ProvideOption{dig.As(new(usecases.DeployCanaryCase))}},
		{constructor: usecases.NewDeployRoutingInteractor, options: []dig.ProvideOption{dig.As(new(usecases.DeployRoutingCase))}},
		{constructor: usecases.NewCodeRouteInteractor, options: []dig.ProvideOption{dig.As(new(usecases.CodeRouteCase))}},
		{constructor: usecases.NewCodeInitInteractor, options: []dig.ProvideOption{dig.As(new(usecases.CodeInitCase))}},
		{constructor: usecases.NewCodeWorkerInteractor, options: []dig.ProvideOption{dig.As(new(usecases.CodeWorkerCase))}},
		{constructor: usecases.NewCodeJobInteractor, options: []dig.ProvideOption{dig.As(new(usecases.CodeJobCase))}},
		{constructor: usecases.NewCodeCronInteractor, options: []dig.ProvideOption{dig.As(new(usecases.CodeCronCase))}},
	}

	for _, p := range providers {
		if err := container.Provide(p.constructor, p.options...); err != nil {
			return nil, err
		}
	}

	return container, nil
}

func resolve[T any](container *dig.Container) (T, error) {
	var value T
	err := container.Invoke(func(resolved T) {
		value = resolved
	})
	return value, err
}
