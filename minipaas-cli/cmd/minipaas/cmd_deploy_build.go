package main

import (
	"github.com/yunier-rojas/minipaas/minipaas-cli/internal/frameworks/logger"
	"github.com/yunier-rojas/minipaas/minipaas-cli/internal/runtime"
)

type DeployBuildArgs struct {
	BaseArgs
}

func (args *DeployBuildArgs) Run() {
	rt, err := runtime.NewDeployBuildRuntime()
	if err != nil {
		logger.PanicErr("Failed to create deploy build runtime", err)
		return
	}

	if err = rt.UseCase.Build(args.Env, args.Verbose); err != nil {
		logger.PanicErr("Failed to build services", err)
		return
	}
}
