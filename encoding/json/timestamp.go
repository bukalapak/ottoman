package json

import (
	"strconv"
	"strings"
	"time"
)

const variousDate = "2006-01-02 15:04:05 -0700"

type Timestamp struct {
	time.Time
}

func (t *Timestamp) UnmarshalJSON(b []byte) error {
	s := string(b)

	if s == `""` || s == "null" {
		return nil
	}

	if strings.Contains(s, " ") {
		if c, err := strconv.Unquote(s); err == nil {
			if v, err := time.Parse(variousDate, c); err == nil {
				t.Time = v
				return nil
			}
		}
	}

	if n, err := strconv.ParseInt(s, 10, 64); err == nil {
		t.Time = time.Unix(n, 0)
		return nil
	}

	return t.Time.UnmarshalJSON(b)
}

func (t *Timestamp) MarshalJSON() ([]byte, error) {
	return t.Time.MarshalJSON()
}
