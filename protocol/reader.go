package protocol

import (
	"encoding/binary"
	"io"
)

// ReadError indicates an error occurred during a read
type ReadError struct {
	error
	msg string
}

func (e ReadError) Error() string {
	return e.msg
}

// ErrorShortRead indicates a failure to read a payload's
// prefix or payload data
var ErrorShortRead = ReadError{
	msg: "Short read",
}

// ErrorEmptyRead indicates a payload size of 0 was read
var ErrorEmptyRead = ReadError{
	msg: "Payload prefix is 0",
}

// Reader handles the read side of the Circuit protocol.
// Specifically, it reads & decodes the payload prefix
// followed by the payload itself.
type Reader struct {
	reader io.Reader
}

// WrapReader wraps a io.Reader with a protocol.Reader
func WrapReader(reader io.Reader) Reader {
	return Reader{
		reader: reader,
	}
}

func (r Reader) Read() ([]byte, error) {
	prefix := make([]byte, 4)
	err := r.readAll(prefix)
	if err != nil {
		return prefix, err
	}
	decoded := decodeUint32(prefix)
	if decoded == 0 {
		return prefix, ErrorEmptyRead
	}
	payload := make([]byte, decoded)
	err = r.readAll(payload)
	return payload, err
}

func (r Reader) readAll(p []byte) error {
	readSize := len(p)
	remaining := readSize
	for tries := 0; tries < 2; tries++ {
		count, err := r.reader.Read(p)
		remaining = remaining - count
		if err != nil {
			if err == io.EOF {
				return ErrorShortRead
			}
			return err
		}
		if remaining == 0 {
			return nil
		}
		p = p[remaining:]
	}
	return ErrorShortRead
}

func decodeUint32(p []byte) uint32 {
	return binary.LittleEndian.Uint32(p)
}
