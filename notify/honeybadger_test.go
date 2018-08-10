package notify_test

import (
	"errors"
	"testing"

	"github.com/bukalapak/ottoman/notify"
	honeybadger "github.com/honeybadger-io/honeybadger-go"
	"github.com/stretchr/testify/assert"
)

func TestHoneybadger(t *testing.T) {
	dc := &NullBackend{}
	hc := honeybadger.New(honeybadger.Configuration{
		Backend: dc,
	})

	h1 := notify.NewHoneybadger(hc, true)
	h2 := notify.NewHoneybadger(hc, false)

	h1.Notify(errors.New("sample error 1"), []byte("<stack>"))
	h2.Notify(errors.New("sample error 2"), []byte("<stack>"))

	hc.Flush()

	assert.Equal(t, []string{"sample error 1"}, dc.Errors)
}

type NullBackend struct {
	Errors []string
}

func (b *NullBackend) Notify(f honeybadger.Feature, p honeybadger.Payload) error {
	b.Errors = append(b.Errors, p.(*honeybadger.Notice).ErrorMessage)
	return nil
}
