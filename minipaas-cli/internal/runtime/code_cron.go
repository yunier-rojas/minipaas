package runtime

import (
	"github.com/yunier-rojas/minipaas/minipaas-cli/internal/core/usecases"
)

type CodeCronRuntime struct {
	UseCase usecases.CodeCronCase
}

func NewCodeCronRuntime() (*CodeCronRuntime, error) {
	container, err := newContainer()
	if err != nil {
		return nil, err
	}

	useCase, err := resolve[usecases.CodeCronCase](container)
	if err != nil {
		return nil, err
	}

	return &CodeCronRuntime{UseCase: useCase}, nil
}
