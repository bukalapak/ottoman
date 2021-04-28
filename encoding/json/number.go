package json

import (
	bjson "encoding/json"
)

type Number struct {
	bjson.Number
}

func (n Number) String() string { return n.Number.String() }

// Float64 returns the number as a float64.
func (n Number) Float64() (float64, error) {
	return n.Number.Float64()
}

// Int64 returns the number as an int64.
func (n Number) Int64() (int64, error) {
	return n.Number.Int64()
}

func (v *Number) UnmarshalJSON(b []byte) error {
	strB := string(b)
	if strB == "" || strB == "null" {
		return nil
	}

	var s bjson.Number
	err := bjson.Unmarshal(b, &s)
	if err != nil {
		return err
	}

	v.Number = s
	return nil
}

func (v *Number) MarshalJSON() ([]byte, error) {
	return []byte(v.String()), nil
}
