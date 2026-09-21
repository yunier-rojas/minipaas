package usecases

type DeployBuildCase interface {
	Build(env string, verbose bool) error
}

type DeployBuildInteractor struct {
	config  ConfigStorePort
	env     EnvironmentPort
	compose ComposePort
	deploy  DeployPort
}

func (i *DeployBuildInteractor) Build(env string, verbose bool) error {
	cfg, err := i.config.Load(env)
	if err != nil {
		return err
	}
	i.env.Apply(env, cfg, verbose)
	return i.deploy.Build(i.compose.FilesForEnv(env, cfg), cfg.Vars["MINIPAAS_IMAGE_PREFIX"], cfg.DeployNamespace(), verbose)
}

func NewDeployBuildInteractor(
	config ConfigStorePort,
	env EnvironmentPort,
	compose ComposePort,
	deploy DeployPort,
) *DeployBuildInteractor {
	return &DeployBuildInteractor{
		config:  config,
		env:     env,
		compose: compose,
		deploy:  deploy,
	}
}

var _ DeployBuildCase = (*DeployBuildInteractor)(nil)
