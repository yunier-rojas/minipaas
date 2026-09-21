package main

import (
	"github.com/yunier-rojas/minipaas/minipaas-cli/internal/frameworks/logger"
	"github.com/yunier-rojas/minipaas/minipaas-cli/internal/runtime"
)

type CodeJobArgs struct {
	BaseArgs
	Services []string `arg:"positional,required" help:"Services to configure as a job."`
}

func (args *CodeJobArgs) Run() {
	rt, err := runtime.NewCodeJobRuntime()
	if err != nil {
		logger.PanicErr("Failed to create code job runtime", err)
		return
	}

	file, err := rt.UseCase.Configure(args.Env, args.Services, args.Verbose)
	if err != nil {
		logger.PanicErr("Failed to configure job", err)
		return
	}

	logger.Info("Updated deploy file: " + file)
}
