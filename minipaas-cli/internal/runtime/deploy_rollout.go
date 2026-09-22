package runtime

import (
	"github.com/yunier-rojas/minipaas/minipaas-cli/internal/core/usecases"
)

type DeployRolloutRuntime struct {
	UseCase usecases.DeployRolloutCase
}

func NewDeployRolloutRuntime() (*DeployRolloutRuntime, error) {
	container, err := newContainer()
	if err != nil {
		return nil, err
	}

	useCase, err := resolve[usecases.DeployRolloutCase](container)
	if err != nil {
		return nil, err
	}

	return &DeployRolloutRuntime{UseCase: useCase}, nil
}
