package json

import (
	"strconv"
	"time"
)

type Timestamp struct {
	time.Time
}

func (t *Timestamp) UnmarshalJSON(b []byte) error {
	if string(b) == `""` || string(b) == "null" {
		return nil
	}

	if n, err := strconv.ParseInt(string(b), 10, 64); err == nil {
		t.Time = time.Unix(n, 0)
		return nil
	}

	return t.Time.UnmarshalJSON(b)
}

func (t *Timestamp) MarshalJSON() ([]byte, error) {
	return t.Time.MarshalJSON()
}
