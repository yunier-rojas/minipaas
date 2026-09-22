package usecases

type CodeJobCase interface {
	Configure(env string, services []string, verbose bool) (string, error)
}

type CodeJobInteractor struct {
	config  ConfigStorePort
	env     EnvironmentPort
	compose ComposePort
}

func (i *CodeJobInteractor) Configure(env string, services []string, verbose bool) (string, error) {
	cfg, err := i.config.Load(env)
	if err != nil {
		return "", err
	}
	i.env.Apply(env, cfg, verbose)
	return i.compose.AddJobDeploy(env, services, cfg.DeployNamespace())
}

func NewCodeJobInteractor(config ConfigStorePort, env EnvironmentPort, compose ComposePort) *CodeJobInteractor {
	return &CodeJobInteractor{config: config, env: env, compose: compose}
}

var _ CodeJobCase = (*CodeJobInteractor)(nil)
