// Package json implements encoding and decoding of JSON.
package json

import (
	"bytes"
	"encoding/json"
	"io"
)

// An Encoder writes JSON values to an output stream.
type Encoder struct {
	*json.Encoder
}

// NewEncoder returns a new encoder that writes to w.
func NewEncoder(w io.Writer) *Encoder {
	c := json.NewEncoder(w)
	c.SetEscapeHTML(false)

	return &Encoder{c}
}

// A Decoder reads and decodes JSON values from an input stream.
type Decoder struct {
	*json.Decoder
}

// NewDecoder returns a new decoder that reads from r.
func NewDecoder(r io.Reader) *Decoder {
	c := json.NewDecoder(r)
	c.UseNumber()

	return &Decoder{c}
}

// Marshal returns the JSON encoding of v.
func Marshal(v interface{}) ([]byte, error) {
	var b bytes.Buffer

	c := NewEncoder(&b)
	err := c.Encode(v)

	return b.Bytes(), err
}

// Unmarshal parses the JSON-encoded data and stores the result in the value pointed to by v.
func Unmarshal(b []byte, v interface{}) error {
	c := NewDecoder(bytes.NewReader(b))
	return c.Decode(v)
}
