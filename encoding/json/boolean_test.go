package json_test

import (
	"strconv"
	"testing"

	"github.com/bukalapak/ottoman/encoding/json"
	"github.com/stretchr/testify/assert"
)

func TestBoolean(t *testing.T) {
	data := map[string]bool{
		`true`:    true,
		`false`:   false,
		`"true"`:  true,
		`"false"`: false,
		`"0"`:     false,
		`"1"`:     true,
		`""`:      false,
		`null`:    false,
	}

	for k, x := range data {
		m := &json.Boolean{}
		err := m.UnmarshalJSON([]byte(k))
		assert.Nil(t, err)
		assert.Equal(t, x, m.Bool())

		z, err := m.MarshalJSON()
		assert.Nil(t, err)

		if x {
			assert.Equal(t, "true", string(z))
		} else {
			assert.Equal(t, "false", string(z))
		}
	}
}

func TestBoolean_numeric(t *testing.T) {
	data := map[int]bool{
		0: false,
		1: true,
	}

	for k, x := range data {
		m := &json.Boolean{}
		err := m.UnmarshalJSON([]byte(strconv.Itoa(k)))
		assert.Nil(t, err)
		assert.Equal(t, x, m.Bool())
	}
}

func TestBoolean_failure(t *testing.T) {
	data := []string{
		`OK`,
		`"ok"`,
		`"wrong"`,
	}

	for _, k := range data {
		m := &json.Boolean{}
		err := m.UnmarshalJSON([]byte(k))
		assert.NotNil(t, err)
		assert.False(t, m.Bool())
	}
}
