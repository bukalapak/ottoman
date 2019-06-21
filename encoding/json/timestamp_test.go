package json_test

import (
	"strconv"
	"testing"
	"time"

	"github.com/bukalapak/ottoman/encoding/json"
	"github.com/stretchr/testify/assert"
)

func TestTimestamp(t *testing.T) {
	testTimestamp(t, []byte("1490318752"))
}

func TestTimestamp_format1(t *testing.T) {
	testTimestamp(t, []byte(`"2017-03-24 08:25:52 +0700"`))
}

func TestTimestamp_fallback(t *testing.T) {
	testTimestamp(t, []byte(`"2017-03-24T08:25:52+07:00"`))
}

func testTimestamp(t *testing.T, b []byte) {
	m := &json.Timestamp{}
	x := "2017-03-24T01:25:52Z"

	err := m.UnmarshalJSON(b)
	assert.Nil(t, err)
	assert.Equal(t, x, m.Time.UTC().Format(time.RFC3339))

	z, err := m.MarshalJSON()
	assert.Nil(t, err)

	s, err := strconv.Unquote(string(z))
	assert.Nil(t, err)

	v, err := time.Parse(time.RFC3339, s)
	assert.Nil(t, err)
	assert.Equal(t, x, v.UTC().Format(time.RFC3339))
}

func TestTimestamp_invalid(t *testing.T) {
	m := &json.Timestamp{}

	err := m.UnmarshalJSON([]byte("BAD-TIME"))
	assert.NotNil(t, err)
	assert.True(t, m.Time.IsZero())
}

func TestTimestamp_emptyQuoteString(t *testing.T) {
	m := &json.Timestamp{}

	err := m.UnmarshalJSON([]byte(`""`))
	assert.Nil(t, err)
	assert.True(t, m.Time.IsZero())
}

func TestTimestamp_null(t *testing.T) {
	m := &json.Timestamp{}

	err := m.UnmarshalJSON([]byte(`null`))
	assert.Nil(t, err)
	assert.True(t, m.Time.IsZero())
}
