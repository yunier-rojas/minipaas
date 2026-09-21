package usecases

type DeployRolloutCase interface {
	Rollout(env string, verbose bool) error
}

type DeployRolloutInteractor struct {
	config  ConfigStorePort
	env     EnvironmentPort
	compose ComposePort
	deploy  DeployPort
}

func (i *DeployRolloutInteractor) Rollout(env string, verbose bool) error {
	cfg, err := i.config.Load(env)
	if err != nil {
		return err
	}
	i.env.Apply(env, cfg, verbose)
	return i.deploy.Rollout(i.compose.FilesForEnv(env, cfg), cfg.DeployNamespace(), verbose)
}

func NewDeployRolloutInteractor(
	config ConfigStorePort,
	env EnvironmentPort,
	compose ComposePort,
	deploy DeployPort,
) *DeployRolloutInteractor {
	return &DeployRolloutInteractor{
		config:  config,
		env:     env,
		compose: compose,
		deploy:  deploy,
	}
}

var _ DeployRolloutCase = (*DeployRolloutInteractor)(nil)
