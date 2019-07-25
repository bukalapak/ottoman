package tracker_test

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/bukalapak/ottoman/tracker"
)

type DummyHTTP struct {
	Resp *http.Response
	Err  error
}

func (dh *DummyHTTP) Do(req *http.Request) (*http.Response, error) {
	return dh.Resp, dh.Err
}

func newDummyHTTP(statusCode int, err error) *DummyHTTP {
	return &DummyHTTP{
		&http.Response{
			StatusCode: statusCode,
			Body:       ioutil.NopCloser(bytes.NewReader([]byte(`{"status": "ok"}`))),
		},
		err,
	}
}

// NewDDCountSeries is function to simplify count creation
func newDDCountSeries(name string, tags []string) *tracker.DDSeries {
	return &tracker.DDSeries{
		Series: []tracker.DDMetric{
			{
				Metric: name,
				Type:   tracker.Count,
				Points: [][2]int64{{time.Now().Unix(), 1}},
				Tags:   tags,
			},
		},
	}
}

func TestDatadog_Track(t *testing.T) {
	type fields struct {
		ServiceName string
		apiKey      string
		httpClient  tracker.HTTPDoer
	}
	type args struct {
		payload interface{}
	}
	tests := map[string]struct {
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		"ok": {
			fields{"testName", "apikey", newDummyHTTP(200, nil)},
			args{newDDCountSeries("name", []string{"tags"})},
			false,
		},
		"ok with nil http": {
			fields{"testName", "apikey", nil},
			args{newDDCountSeries("name", []string{"tags"})},
			false,
		},
		"empty payload": {
			fields{"testName", "apikey", newDummyHTTP(0, tracker.EmptyPayloadErr)},
			args{},
			true,
		},
		"fail request": {
			fields{"testName", "apikey", newDummyHTTP(0, tracker.EmptyPayloadErr)},
			args{newDDCountSeries("name", []string{"tags"})},
			true,
		},
		"bad request": {
			fields{"testName", "apikey", newDummyHTTP(401, nil)},
			args{newDDCountSeries("name", []string{"tags"})},
			true,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			if tt.fields.httpClient == nil {
				tmp := tracker.DefaultHTTPClient
				defer func() {
					tracker.DefaultHTTPClient = tmp
				}()
				tracker.DefaultHTTPClient = newDummyHTTP(200, nil)
			}
			a := tracker.NewDD(
				tt.fields.ServiceName,
				tt.fields.apiKey,
				tt.fields.httpClient,
			)
			if _, err := a.Track(tt.args.payload); (err != nil) != tt.wantErr {
				t.Errorf("Datadog.Track() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
