package clone

import (
	"bytes"
	"io"
	"net/http"
	"net/url"
)

func Request(req *http.Request) *http.Request {
	r := copyRequest(req)
	r.URL = copyURL(req.URL)
	r.Header = copyHeader(req.Header)
	r.Close = false

	return r
}

func URL(u *url.URL) *url.URL {
	return copyURL(u)
}

func DumpBody(b io.ReadCloser) ([]byte, io.ReadCloser, error) {
	z := new(bytes.Buffer)

	if _, err := z.ReadFrom(b); err != nil {
		return nil, b, err
	}

	if err := b.Close(); err != nil {
		return nil, b, err
	}

	return z.Bytes(), nopCloser(z.Bytes()), nil
}

func nopCloser(b []byte) io.ReadCloser {
	return io.NopCloser(bytes.NewReader(b))
}

func copyRequest(req *http.Request) *http.Request {
	r := new(http.Request)
	*r = *req

	return r
}

func copyURL(q *url.URL) *url.URL {
	u := new(url.URL)
	*u = *q

	return u
}

func copyHeader(src http.Header) http.Header {
	h := make(http.Header)
	copyHeaderData(h, src)

	return h
}

func copyHeaderData(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}
