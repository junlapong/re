package runner

import (
	"errors"
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
)

type TRunner struct {
	isKillCalled bool
	killReturn   error

	isStartCalled bool
	startReturn   error
}

func (r *TRunner) KillCommand() error {
	r.isKillCalled = true
	return r.killReturn
}

func (r *TRunner) Start() error {
	r.isStartCalled = true
	return r.startReturn
}

func TestRunnerRun(t *testing.T) {
	t.Run("kill command success then should call Start and return nil", func(t *testing.T) {
		tr := &TRunner{
			killReturn:  nil,
			startReturn: nil,
		}

		err := run(tr)

		assert.Nil(t, err, "should run comamnd success but it have error")
		assert.True(t, tr.isKillCalled, "should have been called Kill command but it not.")
		assert.True(t, tr.isStartCalled, "should have been called Start command but it not.")
	})

	t.Run("should return error when can't start the command", func(t *testing.T) {
		errMsg := "start error"
		tr := &TRunner{
			killReturn:  nil,
			startReturn: errors.New(errMsg),
		}

		err := run(tr)

		assert.Error(t, err, "should return an error but it not.")
	})

	t.Run("should return error when can't kill the command", func(t *testing.T) {
		errMsg := "kill command error"
		tr := &TRunner{
			killReturn:  errors.New(errMsg),
			startReturn: nil,
		}

		err := run(tr)

		assert.Error(t, err, "should return an error but it not.")
	})
}

func TestRunnerStart(t *testing.T) {
	t.Run("should return nil when command execute successfully", func(t *testing.T) {
		task := &Runner{
			prog:   "go",
			args:   []string{"version"},
			stderr: os.Stderr,
			stdout: os.Stdout,
		}

		expectedCmd := exec.Command("go", "version")

		err := task.Start()
		assert.NoError(t, err, "should run comamnd success but it have error")
		assert.Equal(t, expectedCmd.Args, task.cmd.Args, "should run the same command with the initiated one but it doesn't")
	})

	t.Run("should return error when command fail to execute", func(t *testing.T) {
		task := &Runner{
			prog:   "fakecommand",
			args:   []string{"This is not working"},
			stderr: os.Stderr,
			stdout: os.Stdout,
		}

		err := task.Start()
		assert.Error(t, err, "should return an error but it not.")
	})
}

func TestNewRunner(t *testing.T) {
	t.Run("should be Runner instance with prog and args", func(t *testing.T) {
		r := New("go", "version")

		assert.NotNil(t, r, "Runner should be not new")
		assert.Equal(t, "go", r.prog, "should be set prog correctly.")
		assert.Equal(t, []string{"version"}, r.args, "should be set args correctly.")
	})
}

func TestRun(t *testing.T) {
	t.Run("should not be error when call run", func(t *testing.T) {
		r := New("go", "version")

		err := r.Run()

		assert.Nil(t, err, "should not be error.")
	})
}
