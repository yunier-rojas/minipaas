package usecases

type DeployRoutingCase interface {
	Apply(env string, verbose bool) error
}

type DeployRoutingInteractor struct {
	config  ConfigStorePort
	env     EnvironmentPort
	routing RoutingPort
}

func (i *DeployRoutingInteractor) Apply(env string, verbose bool) error {
	cfg, err := i.config.Load(env)
	if err != nil {
		return err
	}
	i.env.Apply(env, cfg, verbose)
	return i.routing.Apply(env, cfg.DeployNamespace(), verbose)
}

func NewDeployRoutingInteractor(
	config ConfigStorePort,
	env EnvironmentPort,
	routing RoutingPort,
) *DeployRoutingInteractor {
	return &DeployRoutingInteractor{
		config:  config,
		env:     env,
		routing: routing,
	}
}

var _ DeployRoutingCase = (*DeployRoutingInteractor)(nil)
