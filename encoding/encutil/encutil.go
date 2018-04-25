package encutil

import (
	"bytes"

	"github.com/bukalapak/ottoman/encoding/json"
	"github.com/bukalapak/ottoman/encoding/msgpack"
)

func MsgPackFromJSON(b []byte) ([]byte, error) {
	var err error
	var buf bytes.Buffer

	var v interface{}

	if err = json.Unmarshal(b, &v); err != nil {
		return nil, err
	}

	enc := msgpack.NewEncoder(&buf)
	err = enc.Encode(v)

	return buf.Bytes(), err
}

func MsgPackToJSON(b []byte) ([]byte, error) {
	var v interface{}

	dec := msgpack.NewDecoder(bytes.NewReader(b))

	if err := dec.Decode(&v); err != nil {
		return nil, err
	}

	return json.Marshal(v)
}
