package protocol

import (
	"bytes"
	"fmt"
	"os/exec"
	"time"
)

type ExecRequest struct {
	Die        bool
	Executable string
	Args       []string
	Options    map[string]string
	Env        map[string]string
	Stdin      []byte
}

// ToExecCommand builds a Go os/exec.Cmd from a source
// ExecRequest
func (er ExecRequest) ToExecCommand() exec.Cmd {
	var command exec.Cmd
	command.Path = er.Executable
	command.Args = er.Args
	command.Env = er.buildEnv()
	command.Stdin = bytes.NewBuffer(er.Stdin)
	return command
}

func (er ExecRequest) buildEnv() []string {
	env := []string{}
	for k, v := range er.Options {
		env = append(env, fmt.Sprintf("%s=%s", k, v))
	}
	for k, v := range er.Env {
		env = append(env, fmt.Sprintf("%s=%s", k, v))
	}
	return env
}

type ExecResult struct {
	Stdout  []byte
	Stderr  []byte
	Success bool
	Elapsed time.Duration
}
