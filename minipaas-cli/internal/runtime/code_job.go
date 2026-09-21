package runtime

import (
	"github.com/yunier-rojas/minipaas/minipaas-cli/internal/core/usecases"
)

type CodeJobRuntime struct {
	UseCase usecases.CodeJobCase
}

func NewCodeJobRuntime() (*CodeJobRuntime, error) {
	container, err := newContainer()
	if err != nil {
		return nil, err
	}

	useCase, err := resolve[usecases.CodeJobCase](container)
	if err != nil {
		return nil, err
	}

	return &CodeJobRuntime{UseCase: useCase}, nil
}
