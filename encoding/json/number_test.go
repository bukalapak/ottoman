package json_test

import (
	"testing"

	"github.com/bukalapak/ottoman/encoding/json"
	"github.com/stretchr/testify/assert"
)

func TestNumber_invalid(t *testing.T) {
	var m json.Number

	x := "lorem"
	err := m.UnmarshalJSON([]byte(x))
	assert.NotNil(t, err)

	x = "1L"
	err = m.UnmarshalJSON([]byte(x))
	assert.NotNil(t, err)

	x = "12.1012.01"
	err = m.UnmarshalJSON([]byte(x))
	assert.NotNil(t, err)
}

func TestNumber_empty_string(t *testing.T) {
	var m json.Number

	x := ""
	err := m.UnmarshalJSON([]byte(x))
	assert.Nil(t, err)
	assert.Equal(t, "", m.String())

	_, err = m.Int64()
	assert.NotNil(t, err)

	b, err := m.MarshalJSON()
	assert.Nil(t, err)
	assert.Equal(t, x, string(b))

	x = "null"
	err = m.UnmarshalJSON([]byte(x))
	assert.Nil(t, err)
	assert.Equal(t, "", m.String())

	_, err = m.Int64()
	assert.Error(t, err)

	b, err = m.MarshalJSON()
	assert.Nil(t, err)
	assert.Equal(t, "", string(b))

	err = m.UnmarshalJSON([]byte(nil))
	assert.Nil(t, err)
	assert.Equal(t, "", m.String())

	_, err = m.Int64()
	assert.Error(t, err)

}

func TestNumber_null_string(t *testing.T) {
	var m json.Number

	x := "null"
	err := m.UnmarshalJSON([]byte(x))
	assert.Nil(t, err)

	_, err = m.Int64()
	assert.Error(t, err)

	b, err := m.MarshalJSON()
	assert.Nil(t, err)
	assert.Equal(t, "", string(b))
}

func TestNumber_number_string(t *testing.T) {
	var m json.Number

	x := "3"
	err := m.UnmarshalJSON([]byte(x))
	assert.Nil(t, err)

	v, err := m.Int64()
	assert.Nil(t, err)
	assert.Equal(t, int64(3), v)

	b, err := m.MarshalJSON()
	assert.Nil(t, err)
	assert.Equal(t, x, string(b))

	x = "1"
	err = m.UnmarshalJSON([]byte(x))
	assert.Nil(t, err)

	v, err = m.Int64()
	assert.Nil(t, err)
	assert.Equal(t, int64(1), v)

	b, err = m.MarshalJSON()
	assert.Nil(t, err)
	assert.Equal(t, x, string(b))

	x = "3.6"
	err = m.UnmarshalJSON([]byte(x))
	assert.Nil(t, err)

	v, err = m.Int64()
	assert.Error(t, err)

	u, err := m.Float64()
	assert.Nil(t, err)
	assert.Equal(t, float64(3.6), u)

	b, err = m.MarshalJSON()
	assert.Nil(t, err)
	assert.Equal(t, x, string(b))
}
