package httputil

import (
	"io"
	"net/http"

	"github.com/bukalapak/ottoman/encoding/encutil"
	"github.com/bukalapak/ottoman/encoding/json"
	"github.com/bukalapak/ottoman/encoding/msgpack"
	"github.com/bukalapak/ottoman/http/header"
)

func DecodeFromHeader(h http.Header, s string, r io.Reader, v interface{}) error {
	if c := header.ContentHeader(h, s); c.Contains(header.MediaTypeMsgPack) {
		return msgpack.NewDecoder(r).Decode(v)
	}

	return json.NewDecoder(r).Decode(v)
}

func EncodeFromHeader(h http.Header, s string, v interface{}) (string, []byte) {
	var t string
	var b []byte

	if z, err := json.Marshal(v); err == nil {
		t = header.MediaTypeJSON
		b = z
	}

	if c, x := EncodeMsgPack(h, s, b); c != "" {
		t = c
		b = x
	}

	return t, b
}

func EncodeMsgPack(h http.Header, s string, b []byte) (string, []byte) {
	if c := header.ContentHeader(h, s); c.Contains(header.MediaTypeMsgPack) {
		if z, err := encutil.MsgPackFromJSON(b); err == nil {
			return c.ContentType().String(), z
		}
	}

	return "", b
}
