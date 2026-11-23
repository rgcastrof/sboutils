package cli

import (
	"fmt"
	"io"
	"os"
	"os/exec"
)

type Executor struct {
	stdout io.Writer
	stderr io.Writer
	cmd *exec.Cmd
}

func NewExecutor() *Executor {
	return &Executor{
		stdout: os.Stdout,
		stderr: os.Stderr,
	}
}

func (e *Executor) Run(cmd string, args []string) error {
	e.cmd = exec.Command(cmd, args...)
	e.cmd.Stdout, e.cmd.Stderr = e.stdout, e.stderr
	if err := e.cmd.Run(); err != nil {
		return fmt.Errorf("%s failed: %w", cmd, err)
	}
	return nil
}
