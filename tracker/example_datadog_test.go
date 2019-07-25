package tracker_test

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/bukalapak/ottoman/tracker"
)

func Example_datadogTracker() {
	tc := tracker.NewDD("service_name",
		"this api key is secret",
		tracker.DDOption{
			Transport: http.DefaultTransport,
			Timeout:   time.Second * 1,
		})
	payload := tracker.DDSeries{
		Series: []tracker.DDMetric{
			{
				Metric: "service_name.command.call",
				Type:   tracker.Count,
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
