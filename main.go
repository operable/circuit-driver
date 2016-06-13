package main

import (
	"bytes"
	"github.com/operable/circuit-driver/protocol"
	"os"
	"os/exec"
	"time"
)

const (
	ExitBadRead  = 2
	ExitBadExec  = 3
	ExitBadWrite = 4
)

func main() {
	decoder := protocol.WrapDecoder(os.Stdin)
	encoder := protocol.WrapEncoder(os.Stdout)
	for {
		var request protocol.ExecRequest
		if err := decoder.Decode(&request); err != nil {
			os.Exit(ExitBadRead)
		}
		execResult, err := executeRequest(request)
		encodeErr := encoder.Encode(execResult)
		if err != nil {
			os.Exit(ExitBadExec)
		}
		if encodeErr != nil {
			os.Exit(ExitBadWrite)
		}

	}
}

func executeRequest(request protocol.ExecRequest) (protocol.ExecResult, error) {
	command := request.ToExecCommand()
	stdout := bytes.NewBuffer([]byte{})
	stderr := bytes.NewBuffer([]byte{})
	command.Stdout = stdout
	command.Stderr = stderr
	start := time.Now()
	err := command.Run()
	finish := time.Now()
	result := protocol.ExecResult{
		Stdout:  stdout.Bytes(),
		Stderr:  stderr.Bytes(),
		Elapsed: finish.Sub(start),
	}
	if err != nil {
		switch execErr := err.(type) {
		case *exec.ExitError:
			result.Success = execErr.Success()
			return result, nil
		default:
			result.Success = false
			return result, err
		}
	}
	return result, nil
}
