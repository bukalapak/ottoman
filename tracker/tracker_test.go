package tracker_test

import (
	"encoding/json"
	"testing"

	"github.com/google/go-cmp/cmp"
)

// dummy is dummy client for testing using tracker without send them to real monitoring service
type dummy struct{}

// Track is implementation Tracker interface to send payload to monitoring server
// for this method just marshal and return result of marshalling
func (d *dummy) Track(payload interface{}) ([]byte, error) {
	b, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// NewNullTracker is a function to create new dummy tracker
func NewNullTracker() *dummy {
	return &dummy{}
}

type dummyPayload struct {
	Name string
}

func TestDummy_Track(t *testing.T) {
	type args struct {
		payload interface{}
	}
	tests := map[string]struct {
		args    args
		want    []byte
		wantErr bool
	}{
		"ok":           {args{dummyPayload{"ottoman"}}, []byte(`{"Name":"ottoman"}`), false},
		"fail marshal": {args{make(chan int)}, nil, true},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			d := NewNullTracker()
			got, err := d.Track(tt.args.payload)
			if (err != nil) != tt.wantErr {
				t.Errorf("Dummy.Track() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("Dummy.Track() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
