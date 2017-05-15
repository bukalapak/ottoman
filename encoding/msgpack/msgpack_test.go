package msgpack_test

import (
	"bytes"
	"testing"

	"github.com/bukalapak/ottoman/encoding/msgpack"
	"github.com/stretchr/testify/assert"
)

var data = []struct {
	m map[string]string
	b []byte
}{
	{map[string]string{}, []byte{0x80}},
	{map[string]string{"hello": "world"}, []byte{0x81, 0xa5, 0x68, 0x65, 0x6c, 0x6c, 0x6f, 0xa5, 0x77, 0x6f, 0x72, 0x6c, 0x64}},
}

func TestDecoder(t *testing.T) {
	for _, v := range data {
		var m map[string]string

		dec := msgpack.NewDecoder(bytes.NewReader(v.b))
		err := dec.Decode(&m)
		assert.Nil(t, err)
		assert.Equal(t, v.m, m)
	}
}

func TestEncoder(t *testing.T) {
	for _, v := range data {
		var b bytes.Buffer

		enc := msgpack.NewEncoder(&b)
		err := enc.Encode(v.m)
		assert.Nil(t, err)
		assert.Equal(t, v.b, b.Bytes())
	}
}

func TestMsgPack(t *testing.T) {
	var b bytes.Buffer
	var m map[string]interface{}

	z := map[string]interface{}{
		"num": uint64(20),
		"foo": "bar",
		"say": map[string]interface{}{
			"hello": "world",
		},
	}

	enc := msgpack.NewEncoder(&b)
	err := enc.Encode(z)
	assert.Nil(t, err)

	dec := msgpack.NewDecoder(&b)
	err = dec.Decode(&m)
	assert.Nil(t, err)
	assert.Equal(t, z, m)
}
