package protocol

import (
	"bytes"
	"encoding/gob"
	"io"
)

// Encoder takes arbitrary Go data, gob encodes it,
// and uses protocol.Writer to write the resulting data
type Encoder struct {
	writer Writer
}

// Encoder configures an Encoder and Writer instance
// around a base io.Writer
func WrapEncoder(w io.Writer) Encoder {
	return Encoder{
		writer: WrapWriter(w),
	}
}

// Encode encodes Go terms and writes the resulting data
func (e Encoder) Encode(p interface{}) error {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(p); err != nil {
		return err
	}
	if _, err := e.writer.Write(buf.Bytes()); err != nil {
		return err
	}
	return nil
}

// Decoder reads data via protocol.Reader, gob decodes the payload,
// and returns the result
type Decoder struct {
	reader Reader
}

// WrapDecoder configures a Decoder and Reader instance
// around a base io.Reader
func WrapDecoder(r io.Reader) Decoder {
	return Decoder{
		reader: WrapReader(r),
	}
}

// Decode reads and decodes data into Go terms
func (d Decoder) Decode(t interface{}) error {
	payload, err := d.reader.Read()
	if err != nil {
		return err
	}
	buf := bytes.NewBuffer(payload)
	dec := gob.NewDecoder(buf)
	if err := dec.Decode(t); err != nil {
		return err
	}
	return nil
}
