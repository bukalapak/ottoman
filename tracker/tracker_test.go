package tracker_test

import (
	"testing"

	"github.com/bukalapak/ottoman/tracker"
	"github.com/google/go-cmp/cmp"
)

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
		// TODO: Add test cases.
		"ok":           {args{dummyPayload{"ottoman"}}, []byte(`{"Name":"ottoman"}`), false},
		"fail marshal": {args{make(chan int)}, nil, true},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			d := &tracker.Dummy{}
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
