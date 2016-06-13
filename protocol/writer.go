package protocol

import (
	"encoding/binary"
	"io"
)

// WriteError indicates an error occurred during a write
type WriteError struct {
	error
	msg string
}

func (we WriteError) Error() string {
	return we.msg
}

// ErrorShortWrite indicates a failure to write a complete payload
var ErrorShortWrite = WriteError{
	msg: "Short write",
}

// Writer handles the write side of the Circuit protocol.
// Specifically, it prefixes all writes with the payload
// size encoded as a little endian 32-bit uint.
type Writer struct {
	writer io.Writer
}

// WrapWriter wraps a io.Writer in a protocol.Writer
func WrapWriter(writer io.Writer) Writer {
	return Writer{
		writer: writer,
	}
}

// Write writes the payload prefix followed by the payload.
// Returns the number of bytes sent and any errors. Note that
// bytes sent includes the payload prefix.
func (w Writer) Write(p []byte) (int, error) {
	prefix := encodeUint32(uint32(len(p)))
	count, err := w.writeAll(prefix)
	if err != nil {
		return count, err
	}
	count, err = w.writeAll(p)
	return count + 4, err
}

func (w Writer) writeAll(p []byte) (int, error) {
	writeSize := len(p)
	offset := 0
	for tries := 0; tries < 2; tries++ {
		count, err := w.writer.Write(p[offset:])
		if err != nil {
			return count, err
		}
		offset = writeSize - (offset + count)
		if offset == 0 {
			return writeSize, nil
		}
	}
	return writeSize - offset, ErrorShortWrite
}

func encodeUint32(i uint32) []byte {
	buf := make([]byte, 4)
	binary.LittleEndian.PutUint32(buf, i)
	return buf
}
