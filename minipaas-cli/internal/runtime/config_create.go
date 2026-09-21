package runtime

import (
	"github.com/yunier-rojas/minipaas/minipaas-cli/internal/core/usecases"
)

type ConfigCreateRuntime struct {
	UseCase usecases.ConfigCreateCase
}

func NewConfigCreateRuntime() (*ConfigCreateRuntime, error) {
	container, err := newContainer()
	if err != nil {
		return nil, err
	}

	useCase, err := resolve[usecases.ConfigCreateCase](container)
	if err != nil {
		return nil, err
	}

	return &ConfigCreateRuntime{UseCase: useCase}, nil
}
