package cli

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
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

func (e *Executor) Run(cmdString string) error {
	split := strings.Split(cmdString, " ")
	e.cmd = exec.Command(split[0], split[1:]...)
	e.cmd.Stdout, e.cmd.Stderr = e.stdout, e.stderr
	if err := e.cmd.Run(); err != nil {
		return fmt.Errorf("%s failed: %w", cmdString, err)
	}
	return nil
}
