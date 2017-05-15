// Package msgpack implements encoding and decoding of MessagePack.
package msgpack

import (
	"io"

	msgpack "gopkg.in/vmihailenco/msgpack.v2"
)

var decodeMapFunc = func(d *msgpack.Decoder) (interface{}, error) {
	n, _ := d.DecodeMapLen()
	m := make(map[string]interface{}, n)

	for i := 0; i < n; i++ {
		k, _ := d.DecodeString()
		v, _ := d.DecodeInterface()

		m[k] = v
	}

	return m, nil
}

// A Decoder reads and decodes MessagePack encoding from an input stream.
type Decoder struct {
	*msgpack.Decoder
}

// NewDecoder returns a new decoder that reads from r.
func NewDecoder(r io.Reader) *Decoder {
	c := msgpack.NewDecoder(r)
	c.DecodeMapFunc = decodeMapFunc

	return &Decoder{c}
}

// Decode reads MessagePack encoding value from its input and stores it in the value pointed to by v.
func (c *Decoder) Decode(v interface{}) error {
	return c.Decoder.Decode(v)
}

// An Encoder writes MessagePack encoding value to an output stream.
type Encoder struct {
	*msgpack.Encoder
}

// NewEncoder returns a new encoder that writes to w.
func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{msgpack.NewEncoder(w)}
}
