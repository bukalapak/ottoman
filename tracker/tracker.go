package tracker

import (
	"errors"
	"net/http"
)

var (
	EmptyPayloadErr = errors.New("empty payload")
	BadRequestErr   = errors.New("bad request")
)

// Tracker is an interface for sending tracking data to server
type Tracker interface {
	Track(data interface{}) ([]byte, error)
}

// HTTPDoer is the interface for http doer to simplify interchange between http client library
type HTTPDoer interface {
	Do(req *http.Request) (resp *http.Response, err error)
}
