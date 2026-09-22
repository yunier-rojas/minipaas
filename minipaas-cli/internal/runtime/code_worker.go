package runtime

import (
	"github.com/yunier-rojas/minipaas/minipaas-cli/internal/core/usecases"
)

type CodeWorkerRuntime struct {
	UseCase usecases.CodeWorkerCase
}

func NewCodeWorkerRuntime() (*CodeWorkerRuntime, error) {
	container, err := newContainer()
	if err != nil {
		return nil, err
	}

	useCase, err := resolve[usecases.CodeWorkerCase](container)
	if err != nil {
		return nil, err
	}

	return &CodeWorkerRuntime{UseCase: useCase}, nil
}
