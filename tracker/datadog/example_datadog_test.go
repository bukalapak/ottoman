package datadog_test

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/bukalapak/ottoman/tracker/datadog"
)

func Example_datadogTracker() {
	tc := datadog.New("service_name",
		"this api key is secret",
		datadog.Option{
			Transport: http.DefaultTransport,
			Timeout:   time.Second * 1,
		})
	payload := datadog.Series{
		Series: []datadog.Metric{
			{
				Metric: "service_name.command.call",
				Type:   datadog.Count,
				Points: [][2]int64{{time.Now().Unix(), 1}},
				Tags: []string{
					"environment:production",
					"version:1",
					"sub_command:mock",
				},
			},
		},
	}
	resp, err := tc.Track(payload)
	if err != nil {
		// do something or just log error
		log.Println(err)
	}
	fmt.Println(resp)
}
