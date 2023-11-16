package influxdb

import (
	"log"
	"strconv"
	"time"
	"trading-engine/consul"
	"trading-engine/errcode"

	influx "github.com/influxdata/influxdb1-client/v2"
)

var influx_cli influx.Client     // 操作influxdb
var PointChan chan *influx.Point //  各个接口通过chan传来的要入influxdb的数据点
var StopChan chan struct{}       // 要关闭PointChan时，用于通知发送端停止发送
var points []*influx.Point       // 暂存一个batch的数据点

const POINTBATCH int = 100000                    // 一次写入influxDB的批大小
const TIME2WRITE time.Duration = time.Second * 3 // 间隔多长时间写入influxDB

func InitInflux() {
	services := consul.DiscoveryService("influxdb-8086", "")
	var addr = services[0].Node.Address + ":" + strconv.Itoa(services[0].Service.Port)
	// log.Printf("%+v", addr)

	var err error
	influx_cli, err = influx.NewHTTPClient(influx.HTTPConfig{
		Addr: "http://" + addr,
	})
	if err != nil {
		log.Printf("Failed to init influxDB: %v", err)
	}
	PointChan = make(chan *influx.Point, 100000)
	StopChan = make(chan struct{})
	points = make([]*influx.Point, 100000)
}

// 将数据点写入influxdb
func writePoints() {
	bp, err := influx.NewBatchPoints(influx.BatchPointsConfig{
		Database: errcode.ServiceName,
	})
	if err != nil {
		log.Printf("Failed to new batchPoints: %v", err)
	}
	// point加入batchPoints，bp写入influxdb
	bp.AddPoints(points)
	err = influx_cli.Write(bp)
	if err != nil {
		log.Printf("Failed to write influxDB: %v", err)
	} else {
		// log.Println("write into influxDB")
	}

	// 清空slice
	points = points[0:0]
}

var timeout chan struct{}

func PointsWriter() {
	// 定时器
	go func() {
		timeout = make(chan struct{})
		// 每过TIME2WRITE秒就把数据写入influxdb
		for {
			time.Sleep(TIME2WRITE)
			timeout <- struct{}{}
		}

	}()

	for {
		select {
		// 定时将数据点写入influxdb
		case <-timeout:
			if len(points) > 0 {
				writePoints()
			}
		case point, ok := <-PointChan:
			if point != nil {
				points = append(points, point)
				if len(points) >= POINTBATCH {
					writePoints()
				}
			} else if !ok {
				if len(points) > 0 {
					writePoints()
				}
				err := influx_cli.Close()
				if err != nil {
					log.Printf("Failed to close influx: %v", err)
				}
				return
			}
		}
	}
}
