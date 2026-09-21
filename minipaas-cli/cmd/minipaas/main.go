package main

import (
	"github.com/alexflint/go-arg"

	"github.com/yunier-rojas/minipaas/minipaas-cli/internal/frameworks/logger"
)

var args struct {
	CodeSubcommand   *CodeSubcommand   `arg:"subcommand:code"`
	DeploySubcommand *DeploySubcommand `arg:"subcommand:deploy"`

	Shell *ShellArgs `arg:"subcommand:shell"`
}

func main() {
	logger.Init()
	arg.MustParse(&args)

	switch {
	case args.CodeSubcommand != nil:
		args.CodeSubcommand.Run()

	case args.DeploySubcommand != nil:
		args.DeploySubcommand.Run()

	case args.Shell != nil:
		args.Shell.Run()
	default:
		logger.Panic("command not supported")
	}
}
