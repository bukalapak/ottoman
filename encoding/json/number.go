package json

import (
	bjson "encoding/json"
	"fmt"
	"strconv"
	"strings"
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
	strB := strings.Trim(string(b), "\"")
	if strB == "" || strB == "null" {
		*v = Number("")
		return nil
	}

	if !isValidNumber(strB) {
		return fmt.Errorf("Invalid number %s for type Number", strB)
	}

	var s bjson.Number
	err := bjson.Unmarshal(b, &s)
	if err != nil {
		return err
	}

	*v = Number(s)
	return nil
}

// isValidNumber reports whether s is a valid JSON number literal. Taken from encoding/json package
func isValidNumber(s string) bool {

	if s == "" {
		return false
	}

	// Optional -
	if s[0] == '-' {
		s = s[1:]
		if s == "" {
			return false
		}
	}

	// Digits
	switch {
	default:
		return false

	case s[0] == '0':
		s = s[1:]

	case '1' <= s[0] && s[0] <= '9':
		s = s[1:]
		for len(s) > 0 && '0' <= s[0] && s[0] <= '9' {
			s = s[1:]
		}
	}

	// . followed by 1 or more digits.
	if len(s) >= 2 && s[0] == '.' && '0' <= s[1] && s[1] <= '9' {
		s = s[2:]
		for len(s) > 0 && '0' <= s[0] && s[0] <= '9' {
			s = s[1:]
		}
	}

	// e or E followed by an optional - or + and
	// 1 or more digits.
	if len(s) >= 2 && (s[0] == 'e' || s[0] == 'E') {
		s = s[1:]
		if s[0] == '+' || s[0] == '-' {
			s = s[1:]
			if s == "" {
				return false
			}
		}
		for len(s) > 0 && '0' <= s[0] && s[0] <= '9' {
			s = s[1:]
		}
	}

	// Make sure we are at the end.
	return s == ""
}
