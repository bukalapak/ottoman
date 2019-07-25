package tracker

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

var (
	DefaultHTTPClient HTTPDoer = &http.Client{}
)

type DDType string

const (
	Count DDType = "count"
	Rate  DDType = "rate"
	Gauge DDType = "gauge"

	ddMetricsURL = "https://api.datadoghq.com/api/v1/series?api_key=%s"
)

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
	httpClient  HTTPDoer
}

// NewDD returns a datadog client with a custom http client
func NewDD(serviceName, apiKey string, httpClient HTTPDoer) *Datadog {

	if httpClient == nil {
		httpClient = DefaultHTTPClient
	}

	return &Datadog{ServiceName: serviceName, apiKey: apiKey, httpClient: httpClient}
}

func (dd *Datadog) Track(payload interface{}) ([]byte, error) {
	if payload == nil {
		return nil, EmptyPayloadErr
	}

	var mapB []byte
	mapB, _ = json.Marshal(payload)
	reader := bytes.NewReader(mapB)

	request, _ := http.NewRequest("POST", fmt.Sprintf(ddMetricsURL, dd.apiKey), reader)

	request.Header.Set("Content-Type", "application/json")

	resp, err := dd.httpClient.Do(request)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode >= 400 {
		return nil, BadRequestErr
	}

	return body, nil
}
