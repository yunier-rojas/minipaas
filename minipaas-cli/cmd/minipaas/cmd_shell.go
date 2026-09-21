package main

import (
	"github.com/yunier-rojas/minipaas/minipaas-cli/internal/frameworks/logger"
	"github.com/yunier-rojas/minipaas/minipaas-cli/internal/runtime"
)

type ShellArgs struct {
	BaseArgs
}

func (args *ShellArgs) Run() {
	rt, err := runtime.NewShellRuntime()
	if err != nil {
		logger.PanicErr("Failed to create shell runtime", err)
		return
	}

	if err = rt.UseCase.Launch(args.Env, args.Verbose); err != nil {
		logger.PanicErr("Failed to launch shell", err)
		return
	}
}
