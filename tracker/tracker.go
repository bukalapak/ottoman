package tracker

import (
	"encoding/json"
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

// Dummy is dummy client for testing using tracker without send them to real monitoring service
type Dummy struct{}

// Track is implementation Tracker interface to send payload to monitoring server
// for this method just marshal and return result of marshalling
func (d *Dummy) Track(payload interface{}) ([]byte, error) {
	b, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	return b, nil
}
