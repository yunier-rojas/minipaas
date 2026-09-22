package main

import (
	"github.com/yunier-rojas/minipaas/minipaas-cli/internal/frameworks/logger"
)

type CodeSubcommand struct {
	CodeInit   *CodeInitArgs   `arg:"subcommand:init"`
	CodeRoute  *CodeRouteArgs  `arg:"subcommand:route"`
	CodeWorker *CodeWorkerArgs `arg:"subcommand:worker"`
	CodeJob    *CodeJobArgs    `arg:"subcommand:job"`
	CodeCron   *CodeCronArgs   `arg:"subcommand:cron"`
}

func (args *CodeSubcommand) Run() {
	switch {
	case args.CodeInit != nil:
		args.CodeInit.Run()
	case args.CodeRoute != nil:
		args.CodeRoute.Run()
	case args.CodeWorker != nil:
		args.CodeWorker.Run()
	case args.CodeJob != nil:
		args.CodeJob.Run()
	case args.CodeCron != nil:
		args.CodeCron.Run()
	default:
		logger.Panic("command not supported")
	}
}
