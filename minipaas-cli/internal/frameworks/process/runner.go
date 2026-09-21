package process

import (
	"bytes"
	"context"
	"os"
	"os/exec"

	"github.com/yunier-rojas/minipaas/minipaas-cli/internal/frameworks/logger"
)

type CommandRunner interface {
	Run(cmd []string, verbose bool) error
	RunWithInput(cmd []string, input []byte, verbose bool) error
	RunOutput(cmd []string, verbose bool) (string, error)
}

type Runner struct{}

func NewRunner() *Runner {
	return &Runner{}
}

func (r *Runner) Run(cmd []string, verbose bool) error {
	if verbose {
		logger.Command(cmd)
	}
	process := exec.CommandContext(context.Background(), cmd[0], cmd[1:]...)
	if verbose {
		process.Stdout = os.Stdout
		process.Stderr = os.Stderr
	}
	return process.Run()
}

func (r *Runner) RunWithInput(cmd []string, input []byte, verbose bool) error {
	if verbose {
		logger.Command(cmd)
	}
	process := exec.CommandContext(context.Background(), cmd[0], cmd[1:]...)
	process.Stdin = bytes.NewReader(input)
	if verbose {
		process.Stdout = os.Stdout
		process.Stderr = os.Stderr
	}
	return process.Run()
}

func (r *Runner) RunOutput(cmd []string, verbose bool) (string, error) {
	if verbose {
		logger.Command(cmd)
	}
	process := exec.CommandContext(context.Background(), cmd[0], cmd[1:]...)
	output, err := process.CombinedOutput()
	return string(output), err
}

var _ CommandRunner = (*Runner)(nil)
