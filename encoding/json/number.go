package json

import (
	"encoding/json"
	bjson "encoding/json"
	"strconv"
)

type Number bjson.Number

// String returns the literal text of the number.
func (n Number) String() string { return string(n) }

// Float64 returns the number as a float64.
func (n Number) Float64() (float64, error) {
	return strconv.ParseFloat(string(n), 64)
}

// Int64 returns the number as an int64.
func (n Number) Int64() (int64, error) {
	return strconv.ParseInt(string(n), 10, 64)
}

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
