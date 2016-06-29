package api

import (
	"bytes"
	"os/exec"
	"time"
)

// Driver is the command execution interface
type Driver interface {
	Run(*ExecRequest) (ExecResult, error)
}

// BlockingDriver executes requests one-at-a-time
type BlockingDriver struct{}

func (bd BlockingDriver) Run(request *ExecRequest) (ExecResult, error) {
	command := request.ToExecCommand()
	stdout := bytes.NewBuffer([]byte{})
	stderr := bytes.NewBuffer([]byte{})
	command.Stdout = stdout
	command.Stderr = stderr
	start := time.Now()
	err := command.Run()
	finish := time.Now()
	result := ExecResult{
		Stdout: stdout.Bytes(),
		Stderr: stderr.Bytes(),
	}
	result.SetElapsed(finish.Sub(start))
	if err != nil {
		switch execErr := err.(type) {
		case *exec.ExitError:
			result.SetSuccess(execErr.Success())
			return result, nil
		default:
			result.SetSuccess(false)
			return result, err
		}
	} else {
		result.SetSuccess(true)
	}
	return result, nil
}
