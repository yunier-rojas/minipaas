package main

import (
	"github.com/yunier-rojas/minipaas/minipaas-cli/internal/frameworks/logger"
	"github.com/yunier-rojas/minipaas/minipaas-cli/internal/runtime"
)

type DeployRoutingArgs struct {
	BaseArgs
}

func (args *DeployRoutingArgs) Run() {
	rt, err := runtime.NewDeployRoutingRuntime()
	if err != nil {
		logger.PanicErr("Failed to create deploy routing runtime", err)
		return
	}

	if err = rt.UseCase.Apply(args.Env, args.Verbose); err != nil {
		logger.PanicErr("Failed to update routing", err)
		return
	}
}
