package runtime

import (
	"github.com/yunier-rojas/minipaas/minipaas-cli/internal/core/usecases"
)

type CodeRouteRuntime struct {
	UseCase usecases.CodeRouteCase
}

func NewCodeRouteRuntime() (*CodeRouteRuntime, error) {
	container, err := newContainer()
	if err != nil {
		return nil, err
	}

	useCase, err := resolve[usecases.CodeRouteCase](container)
	if err != nil {
		return nil, err
	}

	return &CodeRouteRuntime{UseCase: useCase}, nil
}
