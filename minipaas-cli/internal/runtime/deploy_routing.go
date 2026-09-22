package runtime

import (
	"github.com/yunier-rojas/minipaas/minipaas-cli/internal/core/usecases"
)

type DeployRoutingRuntime struct {
	UseCase usecases.DeployRoutingCase
}

func NewDeployRoutingRuntime() (*DeployRoutingRuntime, error) {
	container, err := newContainer()
	if err != nil {
		return nil, err
	}

	useCase, err := resolve[usecases.DeployRoutingCase](container)
	if err != nil {
		return nil, err
	}

	return &DeployRoutingRuntime{UseCase: useCase}, nil
}
