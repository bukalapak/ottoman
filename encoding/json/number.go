package json

import (
	"encoding/json"
	bjson "encoding/json"
)

type Number bjson.Number

func (v *Number) UnmarshalJSON(b []byte) error {
	if string(b) == `""` {
		*v = Number("")
		return nil
	}

	var s bjson.Number
	err := bjson.Unmarshal(b, &s)
	if err != nil {
		return err
	}

	*v = Number(s)
	return nil
}

func (v Number) MarshalJSON() ([]byte, error) {
	n := bjson.Number(v)
	return json.Marshal(n)
}
