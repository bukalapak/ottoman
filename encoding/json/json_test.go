package json_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	jsonx "github.com/bukalapak/ottoman/encoding/json"
	"github.com/stretchr/testify/assert"
)

var data = []struct {
	m map[string]string
	s string
}{
	{map[string]string{}, `{}`},
	{map[string]string{"hello": "world"}, `{"hello":"world"}`},
	{map[string]string{"html": "<p>Awesome!</p>"}, `{"html":"<p>Awesome!</p>"}`},
}

func TestEncoder(t *testing.T) {
	for _, v := range data {
		var b bytes.Buffer

		enc := jsonx.NewEncoder(&b)
		err := enc.Encode(v.m)
		assert.Nil(t, err)
		assert.Equal(t, v.s, strings.TrimSpace(b.String()))
	}
}

func TestDecoder(t *testing.T) {
	for _, v := range data {
		var m map[string]string

		dec := jsonx.NewDecoder(strings.NewReader(v.s))
		err := dec.Decode(&m)
		assert.Nil(t, err)
		assert.Equal(t, v.m, m)
	}

	var z map[string]interface{}

	s := `{"int":10,"float":0.8,"string":"hello"}`
	dec := jsonx.NewDecoder(strings.NewReader(s))
	err := dec.Decode(&z)
	assert.Nil(t, err)
	assert.Equal(t, "hello", z["string"])
	assert.Equal(t, json.Number("0.8"), z["float"])
	assert.Equal(t, json.Number("10"), z["int"])
}

func TestMarshal(t *testing.T) {
	for _, v := range data {
		b, err := jsonx.Marshal(v.m)
		assert.Nil(t, err)
		assert.Equal(t, v.s, strings.TrimSpace(string(b)))
	}
}

func TestUnmarshal(t *testing.T) {
	for _, v := range data {
		var m map[string]string

		err := jsonx.Unmarshal([]byte(v.s), &m)
		assert.Nil(t, err)
		assert.Equal(t, v.m, m)
	}

	var z map[string]interface{}

	s := `{"int":10,"float":0.8,"string":"hello"}`
	err := jsonx.Unmarshal([]byte(s), &z)
	assert.Nil(t, err)
	assert.Equal(t, "hello", z["string"])
	assert.Equal(t, json.Number("0.8"), z["float"])
	assert.Equal(t, json.Number("10"), z["int"])
}
