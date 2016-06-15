package api

import (
	"bytes"
	"testing"
)

func TestEncode(t *testing.T) {
	var buf bytes.Buffer
	epl := WrapEncoder(&buf)
	var td ExecRequest
	td.Executable = "/bin/sayit"
	err := epl.EncodeRequest(&td)
	if err != nil {
		t.Error(err)
	}
	data := buf.Bytes()
	if len(data) < 5 {
		t.Errorf("Unexpected data size: %d", len(data))
	}
}

func TestDecode(t *testing.T) {
	var buf bytes.Buffer
	var requestIn, requestOut ExecRequest
	epl := WrapEncoder(&buf)
	requestIn.Executable = "/bin/sayit"
	requestIn.Env = map[string]interface{}{
		"FOO": 123,
	}
	if err := epl.EncodeRequest(&requestIn); err != nil {
		t.Error(err)
	}
	dpl := WrapDecoder(&buf)
	err := dpl.DecodeRequest(&requestOut)
	if err != nil {
		t.Error(err)
	}
	if requestOut.Executable != requestIn.Executable {
		t.Errorf("Unexpected Executable value: %s", requestOut.Executable)
	}
	if len(requestOut.Env) != 1 {
		t.Errorf("Unexpected env length: %d", len(requestOut.Env))
	}
	if requestOut.Env["FOO"] != 123 {
		t.Errorf("Unexpected payload: %v", requestOut.Env)
	}
}
