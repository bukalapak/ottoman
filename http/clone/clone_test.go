package clone_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync"
	"testing"
	"time"

	httpclone "github.com/bukalapak/ottoman/http/clone"
	"github.com/stretchr/testify/assert"
)

type ctxKey struct {
	Name string
}

func TestRequest(t *testing.T) {
	req := httptest.NewRequest("GET", "https://example.com/search?q=foo", strings.NewReader(`{"foo":"bar"}`))
	req.Header.Set("Content-Type", "application/json")

	key := ctxKey{Name: "hello"}
	val := "world"
	ctx := context.WithValue(req.Context(), key, val)

	req1 := req.WithContext(ctx)
	req2 := httpclone.Request(req1)

	assert.NotEqual(t, fmt.Sprintf("%p", req1), fmt.Sprintf("%p", req2))
	assert.NotEqual(t, fmt.Sprintf("%p", req1.URL), fmt.Sprintf("%p", req2.URL))
	assert.NotEqual(t, fmt.Sprintf("%p", req1.Header), fmt.Sprintf("%p", req2.Header))

	assert.Equal(t, req1.URL.String(), req2.URL.String())
	assert.Equal(t, val, req2.Context().Value(key))
	assert.False(t, req2.Close)
}

func TestRequest_emptyBody(t *testing.T) {
	req1 := httptest.NewRequest("GET", "/", nil)
	req2 := httpclone.Request(req1)

	assert.Equal(t, int64(0), req1.ContentLength)
	assert.Equal(t, int64(0), req2.ContentLength)

	assert.Equal(t, http.NoBody, req1.Body)
	assert.Equal(t, http.NoBody, req2.Body)
}

func TestURL(t *testing.T) {
	u1, _ := url.Parse("https://example.com/search?q=foo")
	u2 := httpclone.URL(u1)

	assert.NotEqual(t, fmt.Sprintf("%p", u1), fmt.Sprintf("%p", u2))
	assert.Equal(t, u1.String(), u2.String())
}

func TestDumpBody(t *testing.T) {
	qs := `{"foo":"bar"}`
	c1 := ioutil.NopCloser(strings.NewReader(qs))
	b2, c2, err := httpclone.DumpBody(c1)
	assert.Nil(t, err)
	assert.Equal(t, qs, string(b2))

	d2 := new(bytes.Buffer)
	d2.ReadFrom(c2)

	assert.Equal(t, qs, d2.String())
	assert.NotEqual(t, fmt.Sprintf("%v", c1), fmt.Sprintf("%v", c2))
}

type BadReader struct{}
type BadCloser struct{}

func (c *BadReader) Close() error               { return nil }
func (c *BadReader) Read(b []byte) (int, error) { return 0, errors.New("read error") }
func (c *BadCloser) Close() error               { return errors.New("close error") }
func (c *BadCloser) Read(b []byte) (int, error) { return 0, io.EOF }

func TestDumpBody_failure(t *testing.T) {
	c1 := &BadReader{}
	b2, c2, err := httpclone.DumpBody(c1)
	assert.Equal(t, "read error", err.Error())
	assert.Nil(t, b2)
	assert.Equal(t, c1, c2)

	c3 := &BadCloser{}
	b4, c4, err := httpclone.DumpBody(c3)
	assert.Equal(t, "close error", err.Error())
	assert.Nil(t, b4)
	assert.Equal(t, c3, c4)
}

func TestRequest_concurrent(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b := new(bytes.Buffer)

		if _, err := b.ReadFrom(r.Body); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		m := map[string]string{"body": b.String()}

		z, err := json.Marshal(m)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Write(z)
	}))

	req, _ := http.NewRequest("GET", srv.URL, strings.NewReader(`{"foo":"bar"}`))
	b, body, err := httpclone.DumpBody(req.Body)
	assert.Nil(t, err)

	req.Body = body // put request body back

	fx := func(t *testing.T, q *http.Response) {
		assert.Equal(t, http.StatusOK, q.StatusCode)
	}

	clientDo(t, req, fx)

	var wg sync.WaitGroup

	for i := 0; i <= 8; i++ {
		wg.Add(1)

		go func(x int) {
			defer wg.Done()

			r := httpclone.Request(req)
			r.Body = httpclone.NopCloser(b)
			r.ContentLength = int64(len(b))

			clientDo(t, r, fx)
		}(i)
	}

	wg.Wait()
	srv.Close()
}

func clientDo(t *testing.T, r *http.Request, fx func(t *testing.T, q *http.Response)) {
	c := http.Client{Timeout: 100 * time.Millisecond}
	b := new(bytes.Buffer)

	q, err := c.Do(r)
	assert.Nil(t, err)

	_, err = b.ReadFrom(q.Body)
	assert.Nil(t, err)

	fx(t, q)
}
