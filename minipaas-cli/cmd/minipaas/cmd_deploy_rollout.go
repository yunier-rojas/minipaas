package main

import (
	"github.com/yunier-rojas/minipaas/minipaas-cli/internal/frameworks/logger"
	"github.com/yunier-rojas/minipaas/minipaas-cli/internal/runtime"
)

type DeployRolloutArgs struct {
	BaseArgs
}

func (args *DeployRolloutArgs) Run() {
	rt, err := runtime.NewDeployRolloutRuntime()
	if err != nil {
		logger.PanicErr("Failed to create deploy rollout runtime", err)
		return
	}

	if err = rt.UseCase.Rollout(args.Env, args.Verbose); err != nil {
		logger.PanicErr("Failed to roll out deployment", err)
		return
	}
}
