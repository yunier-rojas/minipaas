package main

import (
	"github.com/yunier-rojas/minipaas/minipaas-cli/internal/frameworks/logger"
	"github.com/yunier-rojas/minipaas/minipaas-cli/internal/runtime"
)

type DeployConfigArgs struct {
	BaseArgs
	Name string   `arg:"--name" help:"Name (and mounted file name) of the Docker config. Defaults to the base name of the positional file."`
	File string   `arg:"positional" help:"Path to file to use for config content. If omitted, reads from STDIN."`
	For  []string `arg:"--for,separate" help:"Compose services that use the config"`
}

func (args *DeployConfigArgs) Run() {
	rt, err := runtime.NewConfigCreateRuntime()
	if err != nil {
		logger.PanicErr("Failed to create config runtime", err)
		return
	}

	result, err := rt.UseCase.CreateConfig(args.Env, args.Name, args.File, args.For, args.Verbose)
	if err != nil {
		logger.PanicErr("Failed to create config", err)
		return
	}

	if result.Name == "" {
		logger.Warn("Input is empty, skipping.")
		return
	}

	if !result.Attached {
		logger.Warn("The config was created but not attached to any service.")
	}

	logger.Success("Config created: " + result.Name)
}
