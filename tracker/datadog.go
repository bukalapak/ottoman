package tracker

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/pkg/errors"
)

type DDType string

const (
	Count DDType = "count"
	Rate  DDType = "rate"
	Gauge DDType = "gauge"

	ddMetricsURL = "https://api.datadoghq.com/api/v1/series?api_key=%s"
)

// DDOption is the configuration option for the Datadog tracker.
type DDOption struct {
	Transport http.RoundTripper
	Timeout   time.Duration
}

func (do DDOption) httpClient() *http.Client {
	return &http.Client{
		Transport: do.httpTransport(),
		Timeout:   do.httpTimeout(),
	}
}

func (do DDOption) httpTransport() http.RoundTripper {
	if do.Transport == nil {
		return http.DefaultTransport
	}

	return do.Transport
}

func (do DDOption) httpTimeout() time.Duration {
	if do.Timeout == 0 {
		return 3 * time.Second
	}

	return do.Timeout
}

// DatadogMetric represent single metric of series
type DDMetric struct {
	Metric   string     `json:"metric,omitempty" validate:"required"`
	Points   [][2]int64 `json:"points,omitempty" validate:"required"`
	Type     DDType     `json:"type,omitempty"`
	Interval int64      `json:"interval,omitempty"`
	Host     string     `json:"host,omitempty"`
	Tags     []string   `json:"tags,omitempty"`
}

// DDSeries represent time-series data that can be graphed on Datadog’s dashboards
// Further information can be accessed on https://docs.datadoghq.com/api/?lang=bash#post-timeseries-points
type DDSeries struct {
	Series []DDMetric `json:"series"`
}

// Datadog is a client for request to datadog's endpoint
type Datadog struct {
	ServiceName string
	apiKey      string
	option      DDOption
}

// NewDD returns a datadog client with a custom http client
func NewDD(serviceName, apiKey string, option DDOption) *Datadog {
	return &Datadog{ServiceName: serviceName, apiKey: apiKey, option: option}
}

func (dd *Datadog) Track(payload interface{}) ([]byte, error) {
	if payload == nil {
		return nil, EmptyPayloadErr
	}

	var mapB []byte
	mapB, err := json.Marshal(payload)
	if err != nil {
		return nil, errors.Wrap(err, "failed unmarshal")
	}
	reader := bytes.NewReader(mapB)

	request, _ := http.NewRequest("POST", fmt.Sprintf(ddMetricsURL, dd.apiKey), reader)

	request.Header.Set("Content-Type", "application/json")
	resp, err := dd.option.httpClient().Do(request)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		return nil, BadRequestErr
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	return body, nil
}
