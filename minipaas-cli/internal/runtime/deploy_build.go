package runtime

import (
	"github.com/yunier-rojas/minipaas/minipaas-cli/internal/core/usecases"
)

type DeployBuildRuntime struct {
	UseCase usecases.DeployBuildCase
}

func NewDeployBuildRuntime() (*DeployBuildRuntime, error) {
	container, err := newContainer()
	if err != nil {
		return nil, err
	}

	useCase, err := resolve[usecases.DeployBuildCase](container)
	if err != nil {
		return nil, err
	}

	return &DeployBuildRuntime{UseCase: useCase}, nil
}
