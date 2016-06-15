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
	td.Args = []string{"hello", "world"}
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
	requestIn.Args = []string{"how", "now", "brown", "cow"}
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
	if len(requestOut.Args) != 4 {
		t.Errorf("Unexpected args length: %d", len(requestOut.Args))
	}
	if requestOut.Args[0] != "how" || requestOut.Args[1] != "now" ||
		requestOut.Args[2] != "brown" || requestOut.Args[3] != "cow" {
		t.Errorf("Unexpected payload: %v", requestOut.Args)
	}
}
