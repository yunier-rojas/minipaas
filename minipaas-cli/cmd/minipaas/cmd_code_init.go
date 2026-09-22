package main

import (
	"github.com/yunier-rojas/minipaas/minipaas-cli/internal/frameworks/logger"
	"github.com/yunier-rojas/minipaas/minipaas-cli/internal/runtime"
)

type CodeInitArgs struct {
	BaseArgs
	Files     []string `arg:"-c,--compose-file,separate" help:"Extra Compose files to register; defaults to compose.yaml/compose.build.yaml in the project root"`
	Host      string   `arg:"--host" help:"SSH target (user@host) used when --local=false"`
	Local     bool     `arg:"--local" help:"Use the local Docker daemon" default:"true"`
	Registry  string   `arg:"--registry" help:"Optional image prefix for built images, for example ghcr.io/org; when set, deploy build pushes to it"`
	Namespace string   `arg:"--namespace" help:"Deployment namespace used for the stack, services and networks" default:"minipaas"`
}

func (args *CodeInitArgs) Run() {
	rt, err := runtime.NewCodeInitRuntime()
	if err != nil {
		logger.PanicErr("Failed to create code init runtime", err)
		return
	}

	if err = rt.UseCase.Init(args.Env, args.Files, args.Host, args.Local, args.Registry, args.Namespace); err != nil {
		logger.PanicErr("Failed to initialize environment", err)
		return
	}

	logger.Success("Initialized environment: " + args.Env)
}
