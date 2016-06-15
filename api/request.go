package api

import (
	"bytes"
	"fmt"
	"os/exec"
	"time"
)

type ExecRequest struct {
	Die        bool
	Executable string
	Env        map[string]interface{}
	Stdin      []byte
}

// ToExecCommand builds a Go os/exec.Cmd from a source
// ExecRequest
func (er ExecRequest) ToExecCommand() exec.Cmd {
	var command exec.Cmd
	command.Path = er.Executable
	command.Env = er.convertEnv()
	command.Stdin = bytes.NewBuffer(er.Stdin)
	return command
}

func (er ExecRequest) convertEnv() []string {
	retval := []string{}
	for k, v := range er.Env {
		retval = append(retval, fmt.Sprintf("%s=%v", k, v))
	}
	return retval
}

type ExecResult struct {
	Stdout  []byte
	Stderr  []byte
	Success bool
	Elapsed time.Duration
}
