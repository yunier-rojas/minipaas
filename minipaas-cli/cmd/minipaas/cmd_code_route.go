package main

import (
	"github.com/yunier-rojas/minipaas/minipaas-cli/internal/frameworks/logger"
	"github.com/yunier-rojas/minipaas/minipaas-cli/internal/runtime"
)

type CodeRouteArgs struct {
	BaseArgs
	URL    string `arg:"positional,required" help:"Public URL used to expose the service. A path such as /web matches any host and binds plain HTTP on port 80."`
	Target string `arg:"positional,required" help:"Which service to expose. It can also contain the port. Default port to 80."`
}

func (args *CodeRouteArgs) Run() {
	rt, err := runtime.NewCodeRouteRuntime()
	if err != nil {
		logger.PanicErr("Failed to create code route runtime", err)
		return
	}

	routeFile, appsFile, err := rt.UseCase.AddRoute(args.Env, args.URL, args.Target, args.Verbose)
	if err != nil {
		logger.PanicErr("Failed to update routing config", err)
		return
	}

	logger.Info("Updated routing config: " + routeFile)
	logger.Info("Updated deploy file: " + appsFile)
}
