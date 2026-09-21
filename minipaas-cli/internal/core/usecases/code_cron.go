package usecases

type CodeCronCase interface {
	Configure(env string, services []string, cron string, verbose bool) (string, error)
}

type CodeCronInteractor struct {
	config  ConfigStorePort
	env     EnvironmentPort
	compose ComposePort
}

func (i *CodeCronInteractor) Configure(env string, services []string, cron string, verbose bool) (string, error) {
	cfg, err := i.config.Load(env)
	if err != nil {
		return "", err
	}
	i.env.Apply(env, cfg, verbose)
	return i.compose.AddCronDeploy(env, services, cron, cfg.DeployNamespace())
}

func NewCodeCronInteractor(config ConfigStorePort, env EnvironmentPort, compose ComposePort) *CodeCronInteractor {
	return &CodeCronInteractor{config: config, env: env, compose: compose}
}

var _ CodeCronCase = (*CodeCronInteractor)(nil)
