package encutil_test

import (
	"strings"
	"testing"

	"github.com/bukalapak/ottoman/encoding/encutil"
	"github.com/stretchr/testify/assert"
)

var data = map[string][]byte{
	`{}`:                {0x80},
	`{"hello":"world"}`: {0x81, 0xa5, 0x68, 0x65, 0x6c, 0x6c, 0x6f, 0xa5, 0x77, 0x6f, 0x72, 0x6c, 0x64},
	`["foo"]`:           {0x91, 0xa3, 0x66, 0x6f, 0x6f},
}

func TestMsgPackFromJSON(t *testing.T) {
	for s, b := range data {
		v, err := encutil.MsgPackFromJSON([]byte(s))
		assert.Nil(t, err)
		assert.Equal(t, b, v)
	}

	_, err := encutil.MsgPackFromJSON([]byte(`x`))
	assert.NotNil(t, err)
}

func TestMsgPackToJSON(t *testing.T) {
	for s, b := range data {
		v, err := encutil.MsgPackToJSON(b)
		assert.Nil(t, err)
		assert.Equal(t, s, strings.TrimSpace(string(v)))
	}

	_, err := encutil.MsgPackToJSON([]byte{0xa1})
	assert.NotNil(t, err)
}
