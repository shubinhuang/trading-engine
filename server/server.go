package main

import (
	"flag"
	"log"
	"net"
	"net/http"
	_ "net/http/pprof"
	"os"
	"time"
	"trading-engine/consul"
	"trading-engine/engine"
	"trading-engine/errcode"
	"trading-engine/influxdb"
	"trading-engine/process"
	"trading-engine/protos/closeTrade"
	"trading-engine/protos/openTrade"
	"trading-engine/protos/processOrder"
	"trading-engine/rediscache"
	"trading-engine/tradelog"

	influx "github.com/influxdata/influxdb1-client/v2"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
	gops_process "github.com/shirou/gopsutil/process"

	"google.golang.org/grpc"
)

func init() {
	// docker run <image> -servicename=trade0
	flag.StringVar(&errcode.ServiceName, "servicename", "trade0", "service name")
	flag.Parse()

	log.Print("consul Init ...")
	consul.InitConsul()
	log.Print("engine Init ...")
	engine.Init()

	log.Print("redis Init ...")
	rediscache.RedisInit()

	log.Print("logrus Init ...")
	tradelog.InitLogrus()

	log.Print("influx Init ...")
	influxdb.InitInflux()

	log.Print("process Init ...")
	process.Init()
	log.Print("Init done")
}

// 监控server进程的cpu和内存占用，并写入influxdb
func monitor_server() {
	for {
		time.Sleep(time.Second)
		p, _ := gops_process.NewProcess(int32(os.Getpid()))
		//  本进程资源占用率
		cpu_process, _ := p.CPUPercent()
		mem_process, _ := p.MemoryPercent()

		// 系统占用率
		tmp, _ := cpu.Percent(time.Second, false)
		cpu := tmp[0]
		vm, _ := mem.VirtualMemory()
		vm_used := vm.Used
		vm_percent := vm.UsedPercent
		// 往influxdb写入该进程的cpu、内存占用信息
		fields := map[string]interface{}{
			"cpu_process": cpu_process,
			"mem_process": mem_process,
			"cpu":         cpu,
			"vm_used":     int(vm_used),
			"vm_percent":  vm_percent}

		tags := map[string]string{
			"server": errcode.ServiceName, // server 进程
		}
		// 一个数据点（一个请求）
		point, err := influx.NewPoint("monitor", tags, fields, time.Now())
		if err != nil {
			tradelog.Logger.Fatalf("NewPoint error: %v", err)
		}

		select {
		case <-influxdb.StopChan:
			break
		case influxdb.PointChan <- point:
		}
	}
}

func main() {
	// pprof
	go func() {
		tradelog.Logger.Error(http.ListenAndServe("localhost:6060", nil))
	}()

	// 监控server进程的cpu和内存占用，并写入influxdb
	go monitor_server()

	lis, err := net.Listen("tcp", ":5432")
	if err != nil {
		tradelog.Logger.Fatalf("failed to listen: %v", err)
	}

	go influxdb.PointsWriter()

	grpcServer := grpc.NewServer()

	// 服务：开启指定标的的撮合
	openServer := openTrade.Server{}
	// 服务：关闭指定标的的撮合
	closeServer := closeTrade.Server{}
	// 服务：下单/撤单
	orderServer := processOrder.Server{}

	// 将服务注册到grpc中
	openTrade.RegisterOpenServiceServer(grpcServer, &openServer)
	closeTrade.RegisterCloseServiceServer(grpcServer, &closeServer)
	processOrder.RegisterOrderServiceServer(grpcServer, &orderServer)

	log.Printf("gRPC server is listening...")

	// grpc server监听请求
	if err := grpcServer.Serve(lis); err != nil {
		// close(influxdb.StopChan)
		tradelog.Logger.Fatalf("Fail to serve: %v", err)
	}

}
