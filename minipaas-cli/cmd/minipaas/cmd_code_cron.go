package main

import (
	"github.com/yunier-rojas/minipaas/minipaas-cli/internal/frameworks/logger"
	"github.com/yunier-rojas/minipaas/minipaas-cli/internal/runtime"
)

type CodeCronArgs struct {
	BaseArgs
	Services []string `arg:"positional,required" help:"Services to configure as a cron job."`
	Cron     string   `arg:"--cron" help:"Cron schedule." default:"* * * * *"`
}

func (args *CodeCronArgs) Run() {
	rt, err := runtime.NewCodeCronRuntime()
	if err != nil {
		logger.PanicErr("Failed to create code cron runtime", err)
		return
	}

	file, err := rt.UseCase.Configure(args.Env, args.Services, args.Cron, args.Verbose)
	if err != nil {
		logger.PanicErr("Failed to configure cron", err)
		return
	}

	logger.Info("Updated deploy file: " + file)
}
