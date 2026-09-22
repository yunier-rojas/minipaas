package usecases

import "github.com/yunier-rojas/minipaas/minipaas-cli/internal/core/entities"

type CodeRouteCase interface {
	AddRoute(env, url, target string, verbose bool) (routeFile, appsFile string, err error)
}

type CodeRouteInteractor struct {
	config  ConfigStorePort
	env     EnvironmentPort
	routes  RoutePort
	compose ComposePort
}

func (i *CodeRouteInteractor) AddRoute(env, url, target string, verbose bool) (string, string, error) {
	cfg, err := i.config.Load(env)
	if err != nil {
		return "", "", err
	}
	i.env.Apply(env, cfg, verbose)

	routeFile, err := i.routes.AddRoute(env, url, target, cfg.DeployNamespace())
	if err != nil {
		return routeFile, "", err
	}

	service, port := entities.ParseServiceTarget(target)
	appsFile, err := i.compose.AddRouteDeploy(env, service, port, cfg.DeployNamespace())
	if err != nil {
		return routeFile, appsFile, err
	}

	return routeFile, appsFile, nil
}

func NewCodeRouteInteractor(config ConfigStorePort, env EnvironmentPort, routes RoutePort, compose ComposePort) *CodeRouteInteractor {
	return &CodeRouteInteractor{config: config, env: env, routes: routes, compose: compose}
}

var _ CodeRouteCase = (*CodeRouteInteractor)(nil)
