package json

import (
	"encoding/json"
	bjson "encoding/json"
	"strconv"
)

type Number bjson.Number

// String returns the literal text of the number.
func (v Number) String() string { return string(v) }

// Float64 returns the number as a float64.
func (v Number) Float64() (float64, error) {
	return strconv.ParseFloat(string(v), 64)
}

// Int64 returns the number as an int64.
func (v Number) Int64() (int64, error) {
	return strconv.ParseInt(string(v), 10, 64)
}

func (v *Number) UnmarshalJSON(b []byte) error {
	strB := string(b)
	if strB == `""` || strB == "null" {
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

func (v *Number) MarshalJSON() ([]byte, error) {
	if v.String() == "" {
		return json.Marshal(0)
	}

	i, err := v.Int64()
	if err != nil {
		f, err := v.Float64()
		if err != nil {
			return nil, err
		}
		return json.Marshal(f)
	}

	return json.Marshal(i)
}
