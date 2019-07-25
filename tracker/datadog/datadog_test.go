package datadog_test

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/bukalapak/ottoman/tracker"
	"github.com/bukalapak/ottoman/tracker/datadog"
)

type DummyHTTP struct {
	Resp *http.Response
	Err  error
}

func (dh *DummyHTTP) RoundTrip(req *http.Request) (*http.Response, error) {
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

// NewCountSeries is function to simplify count creation
func newCountSeries(name string, tags []string) *datadog.Series {
	return &datadog.Series{
		Series: []datadog.Metric{
			{
				Metric: name,
				Type:   datadog.Count,
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
		option      datadog.Option
	}
	type args struct {
		payload interface{}
	}
	tests := map[string]struct {
		fields  fields
		args    args
		wantErr bool
	}{
		"ok": {
			fields{"testName", "apikey", datadog.Option{Transport: newDummyHTTP(200, nil)}},
			args{newCountSeries("name", []string{"tags"})},
			false,
		},
		"ok with nil http": {
			fields{"testName", "apikey", datadog.Option{Timeout: time.Second}},
			args{newCountSeries("name", []string{"tags"})},
			false,
		},
		"empty payload": {
			fields{"testName", "apikey", datadog.Option{Transport: newDummyHTTP(0, tracker.EmptyPayloadErr)}},
			args{},
			true,
		},
		"fail request": {
			fields{"testName", "apikey", datadog.Option{Transport: newDummyHTTP(0, tracker.EmptyPayloadErr)}},
			args{newCountSeries("name", []string{"tags"})},
			true,
		},
		"bad request": {
			fields{"testName", "apikey", datadog.Option{Transport: newDummyHTTP(401, nil)}},
			args{newCountSeries("name", []string{"tags"})},
			true,
		},
		"bad payload": {
			fields{"testName", "apikey", datadog.Option{Transport: newDummyHTTP(401, nil)}},
			args{make(chan int)},
			true,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			if tt.fields.option.Transport == nil {
				tmp := http.DefaultTransport
				defer func() {
					http.DefaultTransport = tmp
				}()
				http.DefaultTransport = newDummyHTTP(200, nil)
			}

			a := datadog.New(
				tt.fields.ServiceName,
				tt.fields.apiKey,
				tt.fields.option,
			)
			if _, err := a.Track(tt.args.payload); (err != nil) != tt.wantErr {
				t.Errorf("Datadog.Track() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
