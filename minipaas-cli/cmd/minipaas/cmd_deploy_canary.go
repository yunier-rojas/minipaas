package main

import (
	"github.com/yunier-rojas/minipaas/minipaas-cli/internal/frameworks/logger"
	"github.com/yunier-rojas/minipaas/minipaas-cli/internal/runtime"
)

type DeployCanaryArgs struct {
	BaseArgs
	Services []string `arg:"positional,required" help:"Services to release with canary."`
	Replicas int      `arg:"--replicas" help:"Replicas to add for canary release." default:"1"`
}

func (args *DeployCanaryArgs) Run() {
	rt, err := runtime.NewDeployCanaryRuntime()
	if err != nil {
		logger.PanicErr("Failed to create deploy canary runtime", err)
		return
	}

	if err = rt.UseCase.Canary(args.Env, args.Services, args.Replicas, args.Verbose); err != nil {
		logger.PanicErr("Failed to release canary", err)
		return
	}
}
