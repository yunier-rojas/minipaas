package usecases

type CodeWorkerCase interface {
	Configure(env string, services []string, verbose bool) (string, error)
}

type CodeWorkerInteractor struct {
	config  ConfigStorePort
	env     EnvironmentPort
	compose ComposePort
}

func (i *CodeWorkerInteractor) Configure(env string, services []string, verbose bool) (string, error) {
	cfg, err := i.config.Load(env)
	if err != nil {
		return "", err
	}
	i.env.Apply(env, cfg, verbose)
	return i.compose.AddWorkerDeploy(env, services, cfg.DeployNamespace())
}

func NewCodeWorkerInteractor(config ConfigStorePort, env EnvironmentPort, compose ComposePort) *CodeWorkerInteractor {
	return &CodeWorkerInteractor{config: config, env: env, compose: compose}
}

var _ CodeWorkerCase = (*CodeWorkerInteractor)(nil)
