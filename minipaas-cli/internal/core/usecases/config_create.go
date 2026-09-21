package usecases

import "github.com/cockroachdb/errors"

type ConfigCreateCase interface {
	CreateConfig(env, name, file string, services []string, verbose bool) (CreateResult, error)
}

type ConfigCreateInteractor struct {
	config  ConfigStorePort
	env     EnvironmentPort
	input   InputPort
	configs DockerConfigPort
	compose ComposePort
}

func (i *ConfigCreateInteractor) CreateConfig(env, name, file string, services []string, verbose bool) (CreateResult, error) {
	cfg, err := i.config.Load(env)
	if err != nil {
		return CreateResult{}, err
	}
	i.env.Apply(env, cfg, verbose)

	baseName, content, err := i.input.ReadContent(file, name)
	if err != nil {
		return CreateResult{}, err
	}
	if len(content) == 0 {
		return CreateResult{}, nil
	}

	configName := configHashName(baseName, content)

	var svcPerFile map[string][]string
	if len(cfg.Compose.Files) > 0 && len(services) > 0 {
		files := i.compose.FilesForEnv(env, cfg)
		var missing []string
		svcPerFile, missing = i.compose.GroupServicesByFile(files, services, cfg.DeployNamespace())
		if len(missing) > 0 {
			return CreateResult{}, errors.Errorf("services not found: %v", missing)
		}
	}

	if !i.configs.Exists(configName) {
		if err = i.configs.Create(configName, content, verbose); err != nil {
			return CreateResult{}, err
		}
	}

	for f, svcs := range svcPerFile {
		if err = i.compose.AddConfig(f, configName, baseName, svcs, cfg.DeployNamespace()); err != nil {
			return CreateResult{}, errors.Wrapf(err, "update compose file %s", f)
		}
	}

	composeAttached := len(svcPerFile) > 0

	return CreateResult{
		Name:     configName,
		Compose:  composeAttached,
		Attached: composeAttached,
	}, nil
}

func NewConfigCreateInteractor(
	config ConfigStorePort,
	env EnvironmentPort,
	input InputPort,
	configs DockerConfigPort,
	compose ComposePort,
) *ConfigCreateInteractor {
	return &ConfigCreateInteractor{
		config:  config,
		env:     env,
		input:   input,
		configs: configs,
		compose: compose,
	}
}

var _ ConfigCreateCase = (*ConfigCreateInteractor)(nil)
