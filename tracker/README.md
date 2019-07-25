# Tracker

Tracker is a golang library for sending metrics data to the monitoring server

## Usage

```go
import (
    "fmt"
    "github.com/bukalapak/ottoman/tracker"
)

func main(){
   var tc tracker.Tracker
   tc := tracker.NewDD("blcg", "this api key is secret")
   payload := tracker.DDSeries{
        Series: []tracker.DDMetric{
            {
                Metric: "blcg.command.call",
                Type:   tracker.Count,
                Points: [][2]int64{{time.Now().Unix(), 1}},
                Tags:   []string{
                    "environtment:production",
                    "version:1",
                    "command:gen all",
                },
            },
        },
    }
   resp, err:= tc.Track(payload)
   if err!=nil{
       // do something or just log error
    }
    fmt.Println(resp)
}
```
