package usecases

import "github.com/cockroachdb/errors"

type SecretCreateCase interface {
	CreateSecret(env, name, file string, services []string, verbose bool) (CreateResult, error)
}

type SecretCreateInteractor struct {
	config  ConfigStorePort
	env     EnvironmentPort
	input   InputPort
	secrets DockerSecretPort
	compose ComposePort
}

func (i *SecretCreateInteractor) CreateSecret(env, name, file string, services []string, verbose bool) (CreateResult, error) {
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

	secretName := secretHashName(baseName, content)

	var svcPerFile map[string][]string
	if len(cfg.Compose.Files) > 0 && len(services) > 0 {
		files := i.compose.FilesForEnv(env, cfg)
		var missing []string
		svcPerFile, missing = i.compose.GroupServicesByFile(files, services, cfg.DeployNamespace())
		if len(missing) > 0 {
			return CreateResult{}, errors.Errorf("services not found: %v", missing)
		}
	}

	if !i.secrets.Exists(secretName) {
		if err = i.secrets.Create(secretName, content, verbose); err != nil {
			return CreateResult{}, err
		}
	}

	for f, svcs := range svcPerFile {
		if err = i.compose.AddSecret(f, secretName, baseName, svcs, cfg.DeployNamespace()); err != nil {
			return CreateResult{}, errors.Wrapf(err, "update compose file %s", f)
		}
	}

	composeAttached := len(svcPerFile) > 0

	return CreateResult{
		Name:     secretName,
		Compose:  composeAttached,
		Attached: composeAttached,
	}, nil
}

func NewSecretCreateInteractor(
	config ConfigStorePort,
	env EnvironmentPort,
	input InputPort,
	secrets DockerSecretPort,
	compose ComposePort,
) *SecretCreateInteractor {
	return &SecretCreateInteractor{
		config:  config,
		env:     env,
		input:   input,
		secrets: secrets,
		compose: compose,
	}
}

var _ SecretCreateCase = (*SecretCreateInteractor)(nil)
