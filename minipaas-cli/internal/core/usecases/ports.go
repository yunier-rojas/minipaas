package usecases

import "github.com/yunier-rojas/minipaas/minipaas-cli/internal/core/entities"

type ConfigStorePort interface {
	Load(env string) (entities.MinipaasConfig, error)
	Save(env string, cfg entities.MinipaasConfig) (string, error)
}

type ScaffoldPort interface {
	WriteBase(env string, namespace string) (entities.BaseFiles, error)
}

type EnvironmentPort interface {
	Apply(env string, cfg entities.MinipaasConfig, verbose bool)
}

type InputPort interface {
	ReadContent(file, name string) (baseName string, content []byte, err error)
}

type ComposePort interface {
	FilesForEnv(env string, cfg entities.MinipaasConfig) []string
	DiscoverProjectFiles(env string) []string
	GroupServicesByFile(files, services []string, namespace string) (map[string][]string, []string)
	AddConfig(file, configName, baseName string, services []string, namespace string) error
	AddSecret(file, secretName, baseName string, services []string, namespace string) error
	BuildApps(env string, files []string, registry string, namespace string) (string, error)
	AddWorkerDeploy(env string, services []string, namespace string) (string, error)
	AddJobDeploy(env string, services []string, namespace string) (string, error)
	AddCronDeploy(env string, services []string, cron string, namespace string) (string, error)
	AddRouteDeploy(env, service, port string, namespace string) (string, error)
}

type DockerSecretPort interface {
	Exists(name string) bool
	Create(name string, content []byte, verbose bool) error
}

type DockerConfigPort interface {
	Exists(name string) bool
	Create(name string, content []byte, verbose bool) error
}

type ShellPort interface {
	Launch() error
}

type DeployPort interface {
	Build(files []string, registry string, namespace string, verbose bool) error
	Rollout(files []string, namespace string, verbose bool) error
	Canary(files []string, services []string, replicas int, namespace string, verbose bool) error
}

type RoutingPort interface {
	Apply(env string, namespace string, verbose bool) error
}

type RoutePort interface {
	AddRoute(env, url, target string, namespace string) (string, error)
}
