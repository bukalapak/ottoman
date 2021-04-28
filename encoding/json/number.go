package json

import (
	builtin_json "encoding/json"
	"strconv"
)

type Number builtin_json.Number

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
	strB := string(b)
	if strB == "" || strB == "null" {
		return nil
	}

	s := new(builtin_json.Number)
	err := builtin_json.Unmarshal(b, s)
	if err != nil {
		return err
	}

	*v = Number(*s)
	return nil
}

func (v *Number) MarshalJSON() ([]byte, error) {
	return []byte(v.String()), nil
}
