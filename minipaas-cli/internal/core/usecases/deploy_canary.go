package usecases

type DeployCanaryCase interface {
	Canary(env string, services []string, replicas int, verbose bool) error
}

type DeployCanaryInteractor struct {
	config  ConfigStorePort
	env     EnvironmentPort
	compose ComposePort
	deploy  DeployPort
}

func (i *DeployCanaryInteractor) Canary(env string, services []string, replicas int, verbose bool) error {
	cfg, err := i.config.Load(env)
	if err != nil {
		return err
	}
	i.env.Apply(env, cfg, verbose)
	return i.deploy.Canary(i.compose.FilesForEnv(env, cfg), services, replicas, cfg.DeployNamespace(), verbose)
}

func NewDeployCanaryInteractor(
	config ConfigStorePort,
	env EnvironmentPort,
	compose ComposePort,
	deploy DeployPort,
) *DeployCanaryInteractor {
	return &DeployCanaryInteractor{
		config:  config,
		env:     env,
		compose: compose,
		deploy:  deploy,
	}
}

var _ DeployCanaryCase = (*DeployCanaryInteractor)(nil)
