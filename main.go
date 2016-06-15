package main

import (
	"bytes"
	"fmt"
	"github.com/operable/circuit-driver/api"
	"github.com/operable/circuit-driver/util"
	"os"
	"os/exec"
	"time"
)

const (
	ExitBadLogger = 1
	ExitBadRead
	ExitBadExec
	ExitBadWrite
)

func main() {
	inputLogger, err := util.NewDataLogger("/var/log/", util.LogInput, time.Now())
	if err != nil {
		os.Exit(ExitBadLogger)
	}
	outputLogger, err := util.NewDataLogger("/var/log", util.LogOutput, time.Now())
	if err != nil {
		os.Exit(ExitBadLogger)
	}
	decoder := api.WrapDecoder(os.Stdin)
	encoder := api.WrapEncoder(os.Stdout)
	for {
		var request api.ExecRequest
		if err := decoder.DecodeRequest(&request); err != nil {
			inputLogger.WriteString(fmt.Sprintf("Error: %s\n", err))
			//			os.Exit(ExitBadRead)
		}
		inputLogger.WriteString(fmt.Sprintf("request: %+v\n", request))
		if request.Die == true {
			os.Exit(0)
		}
		execResult, err := executeRequest(request)
		outputLogger.WriteString(fmt.Sprintf("result: %+v\n", execResult))
		encodeErr := encoder.EncodeResult(&execResult)
		if err != nil {
			outputLogger.WriteString(fmt.Sprintf("Exec error: %s\n", err))
			//			os.Exit(ExitBadExec)
		}
		if encodeErr != nil {
			outputLogger.WriteString(fmt.Sprintf("Write error: %s\n", err))
			//			os.Exit(ExitBadWrite)
		}

	}
}

func executeRequest(request api.ExecRequest) (api.ExecResult, error) {
	command := request.ToExecCommand()
	stdout := bytes.NewBuffer([]byte{})
	stderr := bytes.NewBuffer([]byte{})
	command.Stdout = stdout
	command.Stderr = stderr
	start := time.Now()
	err := command.Run()
	finish := time.Now()
	result := api.ExecResult{
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
	} else {
		result.Success = true
	}
	return result, nil
}
