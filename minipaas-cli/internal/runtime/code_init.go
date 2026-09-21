package runtime

import (
	"github.com/yunier-rojas/minipaas/minipaas-cli/internal/core/usecases"
)

type CodeInitRuntime struct {
	UseCase usecases.CodeInitCase
}

func NewCodeInitRuntime() (*CodeInitRuntime, error) {
	container, err := newContainer()
	if err != nil {
		return nil, err
	}

	useCase, err := resolve[usecases.CodeInitCase](container)
	if err != nil {
		return nil, err
	}

	return &CodeInitRuntime{UseCase: useCase}, nil
}
