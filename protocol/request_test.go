package protocol

import (
	"bytes"
	"testing"
)

func TestRequestEncode(t *testing.T) {
	var buf bytes.Buffer
	enc := WrapEncoder(&buf)
	dec := WrapDecoder(&buf)
	var request ExecRequest
	request.Executable = "/bin/date"
	var request2 ExecRequest
	enc.Encode(request)
	dec.Decode(&request2)
	if request2.Executable != request.Executable {
		t.Errorf("Expected decoded Executable field to be %s: %s", request.Executable,
			request2.Executable)
	}
	if len(request2.Args) != 0 {
		t.Errorf("Expected args to be empty")
	}
}
