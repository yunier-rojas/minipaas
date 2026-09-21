package runtime

import (
	"github.com/yunier-rojas/minipaas/minipaas-cli/internal/core/usecases"
)

type ShellRuntime struct {
	UseCase usecases.ShellCase
}

func NewShellRuntime() (*ShellRuntime, error) {
	container, err := newContainer()
	if err != nil {
		return nil, err
	}

	useCase, err := resolve[usecases.ShellCase](container)
	if err != nil {
		return nil, err
	}

	return &ShellRuntime{UseCase: useCase}, nil
}
