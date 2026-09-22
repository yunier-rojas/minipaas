package main

import (
	"github.com/yunier-rojas/minipaas/minipaas-cli/internal/frameworks/logger"
	"github.com/yunier-rojas/minipaas/minipaas-cli/internal/runtime"
)

type CodeWorkerArgs struct {
	BaseArgs
	Services []string `arg:"positional,required" help:"Services to configure as a worker."`
}

func (args *CodeWorkerArgs) Run() {
	rt, err := runtime.NewCodeWorkerRuntime()
	if err != nil {
		logger.PanicErr("Failed to create code worker runtime", err)
		return
	}

	file, err := rt.UseCase.Configure(args.Env, args.Services, args.Verbose)
	if err != nil {
		logger.PanicErr("Failed to configure worker", err)
		return
	}

	logger.Info("Updated deploy file: " + file)
}
