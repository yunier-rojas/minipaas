package runtime

import (
	"github.com/yunier-rojas/minipaas/minipaas-cli/internal/core/usecases"
)

type DeployCanaryRuntime struct {
	UseCase usecases.DeployCanaryCase
}

func NewDeployCanaryRuntime() (*DeployCanaryRuntime, error) {
	container, err := newContainer()
	if err != nil {
		return nil, err
	}

	useCase, err := resolve[usecases.DeployCanaryCase](container)
	if err != nil {
		return nil, err
	}

	return &DeployCanaryRuntime{UseCase: useCase}, nil
}
