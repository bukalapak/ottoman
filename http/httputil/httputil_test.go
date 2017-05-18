package httputil_test

import (
	"bytes"
	"net/http"
	"strings"
	"testing"

	"github.com/bukalapak/ottoman/http/httputil"
	"github.com/stretchr/testify/assert"
)

var data = map[string][]byte{
	`{}`:                []byte{0x80},
	`{"hello":"world"}`: []byte{0x81, 0xa5, 0x68, 0x65, 0x6c, 0x6c, 0x6f, 0xa5, 0x77, 0x6f, 0x72, 0x6c, 0x64},
}

func TestEncodeMsgPack(t *testing.T) {
	h := make(http.Header)
	h.Set("Accept", "application/msgpack")

	for s, x := range data {
		c, v := httputil.EncodeMsgPack(h, "Accept", []byte(s))
		assert.Equal(t, "application/msgpack", c)
		assert.Equal(t, x, v)
	}

	h = make(http.Header)
	b := []byte(`{"foo":"bar"}`)
	c, v := httputil.EncodeMsgPack(h, "Accept", b)
	assert.Equal(t, "", c)
	assert.Equal(t, b, v)
}

func TestMsgPackFromJSON(t *testing.T) {
	for s, b := range data {
		v, err := httputil.MsgPackFromJSON([]byte(s))
		assert.Nil(t, err)
		assert.Equal(t, b, v)
	}

	_, err := httputil.MsgPackFromJSON([]byte(`x`))
	assert.NotNil(t, err)
}

func TestMsgPackToJSON(t *testing.T) {
	for s, b := range data {
		v, err := httputil.MsgPackToJSON(b)
		assert.Nil(t, err)
		assert.Equal(t, s, strings.TrimSpace(string(v)))
	}

	_, err := httputil.MsgPackToJSON([]byte{0x0})
	assert.NotNil(t, err)
}

func TestDecodeFromHeader(t *testing.T) {
	for s, b := range data {
		if s == `{}` {
			continue
		}

		m := make(map[string]string)
		h := make(http.Header)
		h.Set("Content-Type", "application/msgpack")

		err := httputil.DecodeFromHeader(h, "Content-Type", bytes.NewBuffer(b), &m)
		assert.Nil(t, err)
		assert.Equal(t, map[string]string{"hello": "world"}, m)
	}

	for s := range data {
		m := make(map[string]string)
		h := make(http.Header)
		h.Set("Content-Type", "application/json")

		err := httputil.DecodeFromHeader(h, "Content-Type", strings.NewReader(s), &m)
		assert.Nil(t, err)

		if s == `{}` {
			assert.Equal(t, map[string]string{}, m)
		} else {
			assert.Equal(t, map[string]string{"hello": "world"}, m)
		}
	}
}

func TestEncodeFromHeader(t *testing.T) {
	x := `{"hello":"world"}`
	m := map[string]string{"hello": "world"}
	h := make(http.Header)
	h.Set("Accept", "application/json")

	c, b := httputil.EncodeFromHeader(h, "Accept", m)
	assert.Equal(t, "application/json", c)
	assert.Equal(t, x, strings.TrimSpace(string(b)))

	h = make(http.Header)
	h.Set("Accept", "application/msgpack")

	c, b = httputil.EncodeFromHeader(h, "Accept", m)
	assert.Equal(t, "application/msgpack", c)
	assert.Equal(t, data[x], b)
}
