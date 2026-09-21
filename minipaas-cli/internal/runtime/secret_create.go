package runtime

import (
	"github.com/yunier-rojas/minipaas/minipaas-cli/internal/core/usecases"
)

type SecretCreateRuntime struct {
	UseCase usecases.SecretCreateCase
}

func NewSecretCreateRuntime() (*SecretCreateRuntime, error) {
	container, err := newContainer()
	if err != nil {
		return nil, err
	}

	useCase, err := resolve[usecases.SecretCreateCase](container)
	if err != nil {
		return nil, err
	}

	return &SecretCreateRuntime{UseCase: useCase}, nil
}
