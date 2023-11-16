package main

import (
	"fmt"
	"strconv"
	pb_close "trading-engine/protos/closeTrade"

	pb_open "trading-engine/protos/openTrade"

	"log"
	"os"
	pb_order "trading-engine/protos/processOrder"

	"github.com/bojand/ghz/printer"
	"github.com/bojand/ghz/runner"
	"github.com/golang/protobuf/proto"
	"github.com/jhump/protoreflect/desc"
)

// 请求数据生成
func dataFuncOpenTrade(mtd *desc.MethodDescriptor, cd *runner.CallData) []byte {
	msg := &pb_open.OpenRequest{}
	msg.Symbol = cd.UUID
	msg.Price = "123.45"

	// rand.Seed(time.Now().UnixNano())
	// r := rand.Intn(100000)
	// msg.Symbol = "s" + strconv.Itoa(r)
	// f := rand.Float64()*1000 + 100
	// msg.Price = strconv.FormatFloat(f, 'f', -1, 64)

	binData, err := proto.Marshal(msg)
	if err != nil {
		log.Fatal(err)
		return []byte{}
	}
	return binData
}

func dataFuncCloseTrade(mtd *desc.MethodDescriptor, cd *runner.CallData) []byte {
	msg := &pb_close.CloseRequest{}
	msg.Symbol = cd.UUID

	binData, err := proto.Marshal(msg)
	if err != nil {
		log.Fatal(err)
		return []byte{}
	}
	return binData
}

// 开启撮合  压测
func ghz_openTrade(c, n, epoch int) {
	// 组装请求数据，编码为BinaryData
	// close_request := pb_close.CloseRequest{}
	// open_request := pb_open.OpenRequest{Symbol: "1", Price: "123.45"}
	// create_request := pb_order.CreateRequest{}
	// cancel_request := pb_order.CancelRequest{}

	// buf := proto.Buffer{}
	// err := buf.EncodeMessage(&open_request)
	// if err != nil {
	// 	log.Fatal(err)
	// 	return
	// }

	report, err := runner.Run(
		// 待测试接口
		"openTrade.OpenService.OpenTrade",
		// 服务器端口
		"localhost:5432",
		// proto文件
		runner.WithProtoFile("../protos/openTrade/openTrade.proto", []string{}),
		// 请求参数
		// runner.WithBinaryData(buf.Bytes()),	// 单个数据
		// runner.WithBinaryDataFunc(dataFuncOpenTrade), // dataFunc生成数据
		// runner.WithDataFromJSON(`[{"symbol":"s1", "price":"123.45"}, {"symbol":"s2", "price":"223.45"},{"symbol":"s3", "price":"323.45"}]`), // 读取json数据
		runner.WithDataFromFile("ghz_opentrade.json"), // 从json文件中读取请求数据
		runner.WithInsecure(true),
		runner.WithTotalRequests(uint(n)),
		// 并发参数
		// runner.WithConcurrencySchedule(runner.ScheduleLine),
		// runner.WithConcurrencyStep(10),
		// runner.WithConcurrencyStart(5),
		// runner.WithConcurrencyEnd(100),
		runner.WithConcurrency(uint(c)),
	)

	if err != nil {
		log.Fatal(err)
		return
	}
	// 指定输出路径
	// file, err := os.Create("report_open_schedule_n100000_1.html")

	file, err := os.Create(fmt.Sprintf("./report/report_open_c%d_n%d_%d.html", c, n, epoch))
	if err != nil {
		log.Fatal(err)
		return
	}
	rp := printer.ReportPrinter{
		Out:    file,
		Report: report,
	}
	// 指定输出格式
	_ = rp.Print("html")
}

// 关闭撮合  压测
func ghz_closeTrade(c, n, epoch int) {

	report, err := runner.Run(
		// 待测试接口
		"closeTrade.CloseService.CloseTrade",
		// 服务器端口
		"localhost:5432",
		// proto文件
		runner.WithProtoFile("../protos/closeTrade/closeTrade.proto", []string{}),
		// 请求参数
		// runner.WithBinaryDataFunc(dataFuncCloseTrade),
		runner.WithDataFromFile("ghz_closetrade.json"),
		runner.WithInsecure(true),
		runner.WithTotalRequests(uint(n)),
		// 并发参数
		// runner.WithConcurrencySchedule(runner.ScheduleLine),
		// runner.WithConcurrencyStep(10),
		// runner.WithConcurrencyStart(5),
		// runner.WithConcurrencyEnd(100),
		runner.WithConcurrency(uint(c)),
	)

	if err != nil {
		log.Fatal(err)
		return
	}
	// 指定输出路径
	file, err := os.Create(fmt.Sprintf("./report/report_close_c%d_n%d_%d.html", c, n, epoch))
	if err != nil {
		log.Fatal(err)
		return
	}
	rp := printer.ReportPrinter{
		Out:    file,
		Report: report,
	}
	// 指定输出格式
	_ = rp.Print("html")
}

// 下单  压测
func ghz_createOrder(c, n, epoch int, datafile string) {

	create_request := pb_order.CreateRequest{}
	create_request.OrderId = "1"

	report, err := runner.Run(
		// 待测试接口
		"processOrder.OrderService.CreateOrder",
		// 服务器端口
		"localhost:5432",
		// proto文件
		runner.WithProtoFile("../protos/processOrder/processOrder.proto", []string{}),
		// 请求参数
		runner.WithDataFromFile(datafile),
		runner.WithInsecure(true),
		runner.WithTotalRequests(uint(n)),
		// 并发参数
		// runner.WithConcurrencySchedule(runner.ScheduleLine),
		// runner.WithConcurrencyStep(10),
		// runner.WithConcurrencyStart(5),
		// runner.WithConcurrencyEnd(100),
		runner.WithConcurrency(uint(c)),
	)

	if err != nil {
		log.Fatal(err)
		return
	}
	// 指定输出路径
	file, err := os.Create(fmt.Sprintf("./report/report_create_c%d_n%d_%d.html", c, n, epoch))
	if err != nil {
		log.Fatal(err)
		return
	}
	rp := printer.ReportPrinter{
		Out:    file,
		Report: report,
	}
	// 指定输出格式
	_ = rp.Print("html")
}

// 撤单  压测
func ghz_cancelOrder(c, n, epoch int) {
	cancel_request := pb_order.CancelRequest{}
	cancel_request.OrderId = "1"

	report, err := runner.Run(
		// 待测试接口
		"processOrder.OrderService.CancelOrder",
		// 服务器端口
		"localhost:5432",
		// proto文件
		runner.WithProtoFile("../protos/processOrder/processOrder.proto", []string{}),
		// 请求参数
		runner.WithDataFromFile("ghz_cancelorder.json"),
		runner.WithInsecure(true),
		runner.WithTotalRequests(uint(n)),
		// 并发参数
		// runner.WithConcurrencySchedule(runner.ScheduleLine),
		// runner.WithConcurrencyStep(10),
		// runner.WithConcurrencyStart(5),
		// runner.WithConcurrencyEnd(100),
		runner.WithConcurrency(uint(c)),
	)

	if err != nil {
		log.Fatal(err)
		return
	}
	// 指定输出路径
	file, err := os.Create(fmt.Sprintf("./report/report_cancel_c%d_n%d_%d.html", c, n, epoch))
	if err != nil {
		log.Fatal(err)
		return
	}
	rp := printer.ReportPrinter{
		Out:    file,
		Report: report,
	}
	// 指定输出格式
	_ = rp.Print("html")
}

func main() {
	if len(os.Args) >= 5 {
		c, err := strconv.Atoi(os.Args[2])
		if err != nil {
			return
		}
		n, err := strconv.Atoi(os.Args[3])
		if err != nil {
			return
		}
		epoch, err := strconv.Atoi(os.Args[4])
		if err != nil {
			return
		}

		if os.Args[1] == "opentrade" {
			ghz_openTrade(c, n, epoch)
		} else if os.Args[1] == "closetrade" {
			ghz_closeTrade(c, n, epoch)
		} else if os.Args[1] == "createorder" {
			ghz_createOrder(c, n, epoch, os.Args[5])
		} else if os.Args[1] == "cancelorder" {
			ghz_cancelOrder(c, n, epoch)
		} else {
			fmt.Println("wrong command")
		}
	} else {
		fmt.Println("wrong command")
	}
}
