package notify

import (
	honeybadger "github.com/honeybadger-io/honeybadger-go"
)

// Honeybadger wraps honeybadger.Client
type Honeybadger struct {
	client  *honeybadger.Client
	enabled bool
}

func NewHoneybadger(client *honeybadger.Client, enabled bool) *Honeybadger {
	return &Honeybadger{client: client, enabled: enabled}
}

// Notify reports the err to the Honeybadger service.
func (h *Honeybadger) Notify(err interface{}, stack []byte) {
	if h.enabled {
		h.client.Notify(err, stack)
	}
}
