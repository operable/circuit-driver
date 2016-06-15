package api

import (
	"encoding/gob"
	circuit "github.com/operable/circuit-driver/io"
	"io"
)

// Encoder takes arbitrary Go data, gob encodes it,
// and uses protocol.Writer to write the resulting data
type Encoder struct {
	encoder *gob.Encoder
}

// Encoder configures an Encoder and Writer instance
// around a base io.Writer
func WrapEncoder(w io.Writer) Encoder {
	return Encoder{
		encoder: gob.NewEncoder(circuit.NewCircuitWriter(w)),
	}
}

// EncodeRequest encodes ExecRequests and writes them to the underlying
// transport via protocol.Writer
func (e Encoder) EncodeRequest(request *ExecRequest) error {
	return e.encode(request)
}

// EncodeResult encodes ExecResults and writes them to the underlying
// transport via protcol.Writer
func (e Encoder) EncodeResult(result *ExecResult) error {
	return e.encode(result)
}

func (e Encoder) encode(p interface{}) error {
	return e.encoder.Encode(p)
}

// Decoder reads data via protocol.Reader, gob decodes the payload,
// and returns the result
type Decoder struct {
	decoder *gob.Decoder
}

// WrapDecoder configures a Decoder and Reader instance
// around a base io.Reader
func WrapDecoder(r io.Reader) Decoder {
	return Decoder{
		decoder: gob.NewDecoder(circuit.NewCircuitReader(r)),
	}
}

// DecodeRequest reads and decodes ExecRequests
func (d Decoder) DecodeRequest(request *ExecRequest) error {
	return d.decode(request)
}

// DecodeResult reads and decodes ExecResults
func (d Decoder) DecodeResult(result *ExecResult) error {
	return d.decode(result)
}

func (d Decoder) decode(p interface{}) error {
	return d.decoder.Decode(p)
}
