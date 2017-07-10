package encutil

import (
	"bytes"

	"github.com/bukalapak/ottoman/encoding/json"
	"github.com/bukalapak/ottoman/encoding/msgpack"
)

func MsgPackFromJSON(b []byte) ([]byte, error) {
	var err error
	var buf bytes.Buffer

	var m map[string]interface{}

	if err = json.Unmarshal(b, &m); err != nil {
		return nil, err
	}

	enc := msgpack.NewEncoder(&buf)
	err = enc.Encode(m)

	return buf.Bytes(), err
}

func MsgPackToJSON(b []byte) ([]byte, error) {
	var m map[string]interface{}

	dec := msgpack.NewDecoder(bytes.NewReader(b))

	if err := dec.Decode(&m); err != nil {
		return nil, err
	}

	return json.Marshal(m)
}
