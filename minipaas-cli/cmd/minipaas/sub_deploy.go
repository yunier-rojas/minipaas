package main

import (
	"github.com/yunier-rojas/minipaas/minipaas-cli/internal/frameworks/logger"
)

type DeploySubcommand struct {
	DeployBuild   *DeployBuildArgs   `arg:"subcommand:build"`
	DeployRollout *DeployRolloutArgs `arg:"subcommand:rollout"`
	DeployCanary  *DeployCanaryArgs  `arg:"subcommand:canary"`
	DeployRouting *DeployRoutingArgs `arg:"subcommand:routing"`
	DeploySecret  *DeploySecretArgs  `arg:"subcommand:secret"`
	DeployConfig  *DeployConfigArgs  `arg:"subcommand:config"`
}

func (args *DeploySubcommand) Run() {
	switch {
	case args.DeployBuild != nil:
		args.DeployBuild.Run()
	case args.DeployRollout != nil:
		args.DeployRollout.Run()
	case args.DeployCanary != nil:
		args.DeployCanary.Run()
	case args.DeployRouting != nil:
		args.DeployRouting.Run()
	case args.DeploySecret != nil:
		args.DeploySecret.Run()
	case args.DeployConfig != nil:
		args.DeployConfig.Run()
	default:
		logger.Panic("command not supported")
	}
}
