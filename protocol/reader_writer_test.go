package protocol

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"math/rand"
	"testing"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type badWriter struct {
	writer io.Writer
}

type badReader struct {
	reader io.Reader
}

func (bw badWriter) Write(p []byte) (int, error) {
	offset := int(math.Min(0.0, float64(len(p)-2)))
	if offset > 0 {
		offset = rand.Intn(offset - 1)
		fmt.Printf("Offset: %d\np length: %d\n", offset, len(p))
		return bw.writer.Write(p[:offset])
	}
	return offset, nil
}

func (br badReader) Read(p []byte) (int, error) {
	count, err := br.reader.Read(p)
	if err == io.EOF {
		return 0, nil
	}
	return count, err
}

func TestReadWrite(t *testing.T) {
	payload := []byte("hello")
	expected := len(payload) + 4
	var buf bytes.Buffer
	protocolReader := WrapReader(&buf)
	protocolWriter := WrapWriter(&buf)
	count, err := protocolWriter.Write(payload)
	if err != nil {
		t.Error(err)
	}
	if count != expected {
		t.Errorf("Expected to write %d bytes; Wrote %d bytes", expected, count)
	}
	payload2, err := protocolReader.Read()
	if err != nil {
		t.Error(err)
	}
	if len(payload2) != (expected - 4) {
		t.Errorf("Expected to read %d bytes; Read %d bytes", expected, len(payload2))
	}
	compareData(payload2, payload, t)
}

func TestReadNoData(t *testing.T) {
	var buf bytes.Buffer
	protocolReader := WrapReader(&buf)
	_, err := protocolReader.Read()
	if err != ErrorShortRead {
		t.Error(err)
	}
}

func TestReadShortPrefix(t *testing.T) {
	var buf bytes.Buffer
	prefix := make([]byte, 4)
	binary.LittleEndian.PutUint32(prefix, 15)
	buf.Write(prefix[0:1])
	protocolReader := WrapReader(&buf)
	_, err := protocolReader.Read()
	if err != ErrorShortRead {
		t.Error(err)
	}
}

func TestReadShortPayload(t *testing.T) {
	payload := []byte("hello")
	var buf bytes.Buffer
	prefix := make([]byte, 4)
	binary.LittleEndian.PutUint32(prefix, uint32(len(payload)))
	buf.Write(prefix)
	buf.Write(payload[0:3])
	protocolReader := WrapReader(&buf)
	readBuf, err := protocolReader.Read()
	if err != ErrorShortRead {
		t.Error(err)
	}
	compareData(readBuf[0:3], payload, t)
}

func TestShortWrites(t *testing.T) {
	payload := []byte("hello")
	var buf bytes.Buffer
	bw := badWriter{
		writer: &buf,
	}
	protocolWriter := WrapWriter(bw)
	count, err := protocolWriter.Write(payload)
	if err != ErrorShortWrite {
		t.Errorf("Expected ErrorShortWrite")
	}
	if count == len(payload) {
		t.Errorf("Expected count != payload size")
	}
}

func compareData(readData, writeData []byte, t *testing.T) {
	for i, v := range readData {
		if v != writeData[i] {
			t.Errorf("Read vs. write data differs at offset %d: %d %d", i, v, writeData[i])
		}
	}
}

func TestReadBadPrefix(t *testing.T) {
	payload := []byte("hello")
	badLength := uint32(len(payload) + 1)
	prefix := make([]byte, 4)
	binary.LittleEndian.PutUint32(prefix, badLength)
	var buf bytes.Buffer
	buf.Write(prefix)
	buf.Write(payload)
	protocolReader := WrapReader(&buf)
	_, err := protocolReader.Read()
	if err != ErrorShortRead {
		t.Errorf("Expecting ErrorShortRead: %s", err)
	}
}

func TestFlakyReader(t *testing.T) {
	payload := []byte("hello")
	badLength := uint32(len(payload) + 1)
	prefix := make([]byte, 4)
	binary.LittleEndian.PutUint32(prefix, badLength)
	var buf bytes.Buffer
	buf.Write(prefix)
	buf.Write(payload)
	protocolReader := WrapReader(badReader{
		reader: &buf,
	})
	_, err := protocolReader.Read()
	if err != ErrorShortRead {
		t.Errorf("Expecting ErrorShortRead: %s", err)
	}
}
