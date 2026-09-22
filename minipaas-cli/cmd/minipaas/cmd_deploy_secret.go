package main

import (
	"github.com/yunier-rojas/minipaas/minipaas-cli/internal/frameworks/logger"
	"github.com/yunier-rojas/minipaas/minipaas-cli/internal/runtime"
)

type DeploySecretArgs struct {
	BaseArgs
	Name string   `arg:"--name" help:"Name (and mounted file name) of the Docker secret. Defaults to the base name of the positional file."`
	File string   `arg:"positional" help:"Path to file to use for secret content. If omitted, reads from STDIN."`
	For  []string `arg:"--for,separate" help:"Compose services that use the secret"`
}

func (args *DeploySecretArgs) Run() {
	rt, err := runtime.NewSecretCreateRuntime()
	if err != nil {
		logger.PanicErr("Failed to create secret runtime", err)
		return
	}

	result, err := rt.UseCase.CreateSecret(args.Env, args.Name, args.File, args.For, args.Verbose)
	if err != nil {
		logger.PanicErr("Failed to create secret", err)
		return
	}

	if result.Name == "" {
		logger.Warn("Input is empty, skipping.")
		return
	}

	if !result.Attached {
		logger.Warn("The secret was created but not attached to any service.")
	}

	logger.Success("Secret created: " + result.Name)
}
