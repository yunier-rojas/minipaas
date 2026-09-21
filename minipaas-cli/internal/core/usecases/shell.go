package usecases

type ShellCase interface {
	Launch(env string, verbose bool) error
}

type ShellInteractor struct {
	config ConfigStorePort
	env    EnvironmentPort
	shell  ShellPort
}

func (i *ShellInteractor) Launch(env string, verbose bool) error {
	cfg, err := i.config.Load(env)
	if err != nil {
		return err
	}
	i.env.Apply(env, cfg, verbose)
	return i.shell.Launch()
}

func NewShellInteractor(config ConfigStorePort, env EnvironmentPort, shell ShellPort) *ShellInteractor {
	return &ShellInteractor{config: config, env: env, shell: shell}
}

var _ ShellCase = (*ShellInteractor)(nil)
