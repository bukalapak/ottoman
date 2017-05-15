package json

import (
	"encoding/json"
	"strconv"
)

type Boolean struct {
	b bool
}

func (v *Boolean) UnmarshalJSON(b []byte) error {
	if string(b) == `""` || string(b) == "null" {
		return nil
	}

	var err error

	s, _ := strconv.Unquote(string(b))

	if v.b, err = strconv.ParseBool(s); err != nil {
		if v.b, err = strconv.ParseBool(string(b)); err != nil {
			return err
		}
	}

	return nil
}

func (v Boolean) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.b)
}

func (v Boolean) Bool() bool {
	return v.b
}
