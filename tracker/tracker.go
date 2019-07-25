// package tracker implements sending tracking data to monitoring server
package tracker

import (
	"errors"
)

var (
	EmptyPayloadErr = errors.New("empty payload")
	BadRequestErr   = errors.New("bad request")
)

// Tracker is an interface for sending tracking data to the monitoring server
type Tracker interface {
	Track(data interface{}) ([]byte, error)
}
