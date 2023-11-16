package main

import (
	"log"
	"time"
	"trading-engine/influxdb"

	influx "github.com/influxdata/influxdb1-client/v2"
)

func writeP() {
	for i := 1; i < 25; i++ {
		time.Sleep(time.Microsecond)
		tags := map[string]string{
			"cpu":    "ih-cpu",
			"server": "test",
		}
		fields := map[string]interface{}{
			"idle":   201.1,
			"system": 43.3,
			"user":   i,
		}
		pt, err := influx.NewPoint("cpu_usage", tags, fields, time.Now())
		if err != nil {
			log.Fatalf("NewPoint error: %v", err)
		}
		select {
		case <-influxdb.StopChan:
			break
		case influxdb.PointChan <- pt:
		}

	}
}

// func main() {

// 	influxdb.InitInflux()
// 	go influxdb.PointsWriter()
// 	writeP()
// 	time.Sleep(time.Second * 2)
// }
