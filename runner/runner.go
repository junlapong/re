package runner

import (
	"errors"
	"io"
	"os"
	"os/exec"
	"syscall"
	"time"
)

// Runner is the task runnde
type Runner struct {
	prog   string
	args   []string
	cmd    *exec.Cmd
	stdout io.Writer
	stderr io.Writer
}

var task *Runner

// New creates new task runner if not exists
func New(prog string, args ...string) *Runner {
	if task == nil {
		return &Runner{
			prog:   prog,
			args:   args,
			stderr: os.Stderr,
			stdout: os.Stdout,
		}
	}

	return task
}

func (r *Runner) Start() error {
	cmd := exec.Command(r.prog, r.args...)
	cmd.Stdout = r.stdout
	cmd.Stderr = r.stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	r.cmd = cmd

	return r.cmd.Start()
}

type iRunner interface {
	Start() error
	KillCommand() error
}

// Run starts the runner
func (r *Runner) Run() error {
	return run(r)
}

func run(r iRunner) error {
	err := r.KillCommand()
	if err != nil {
		return err
	}

	err = r.Start()
	if err != nil {
		return err
	}

	return nil
}

func (r *Runner) KillCommand() error {
	if r.cmd == nil {
		return nil
	}

	if r.cmd.Process == nil {
		return nil
	}

	pid := r.cmd.Process.Pid

	done := make(chan struct{})
	go func() {
		r.cmd.Wait()
		close(done)
	}()

	// try soft kill
	syscall.Kill(-pid, syscall.SIGINT)
	select {
	case <-time.After(3 * time.Second):
		// go hard because soft is not always the solution
		err := syscall.Kill(-pid, syscall.SIGKILL)
		if err != nil {
			return errors.New("Fail killing on going process")
		}
	case <-done:
	}

	return nil
}
