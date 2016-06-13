package protocol

import (
	"bytes"
	"testing"
)

type TestData struct {
	Kind    int
	Payload []string
}

func TestEncode(t *testing.T) {
	var buf bytes.Buffer
	epl := WrapEncoder(&buf)
	var td TestData
	td.Kind = 3
	td.Payload = []string{"hello", "world"}
	err := epl.Encode(td)
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
	var tdIn, tdOut TestData
	epl := WrapEncoder(&buf)
	tdIn.Kind = 5
	tdIn.Payload = []string{"how", "now", "brown", "cow"}
	if err := epl.Encode(tdIn); err != nil {
		t.Error(err)
	}
	dpl := WrapDecoder(&buf)
	err := dpl.Decode(&tdOut)
	if err != nil {
		t.Error(err)
	}
	if tdOut.Kind != 5 {
		t.Errorf("Unexpected Kind value: %d", tdOut.Kind)
	}
	if len(tdOut.Payload) != 4 {
		t.Errorf("Unexpected payload length: %d", len(tdOut.Payload))
	}
	if tdOut.Payload[0] != "how" || tdOut.Payload[1] != "now" ||
		tdOut.Payload[2] != "brown" || tdOut.Payload[3] != "cow" {
		t.Errorf("Unexpected payload: %v", tdOut.Payload)
	}
}

func TestBadEncode(t *testing.T) {
	var buf bytes.Buffer
	var tdIn TestData
	epl := Encoder{
		writer: WrapWriter(badWriter{
			writer: &buf,
		}),
	}
	tdIn.Kind = 1
	err := epl.Encode(tdIn)
	if err != ErrorShortWrite {
		t.Error(err)
	}
}

func TestEmptyDecode(t *testing.T) {
	var buf bytes.Buffer
	var tdOut TestData
	dpl := Decoder{
		reader: WrapReader(&buf),
	}
	err := dpl.Decode(&tdOut)
	if err != ErrorShortRead {
		t.Error(err)
	}
}
