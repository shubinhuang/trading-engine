package closeTrade

import (
	"context"
	"strconv"
	"time"
	"trading-engine/errcode"
	"trading-engine/influxdb"
	"trading-engine/process"
	"trading-engine/tradelog"

	influx "github.com/influxdata/influxdb1-client/v2"
)

// 实现服务中定义的CloseTrade方法
type Server struct {
	UnimplementedCloseServiceServer
}

func (server *Server) CloseTrade(ctx context.Context, in *CloseRequest) (*CloseReply, error) {
	time_start := time.Now()
	reply := CloseReply{}

	symbol := in.GetSymbol()
	errorCode := process.CloseEngine(symbol)
	// reply := CloseReply{Code: errorCode.Code, Msg: errorCode.Msg}
	reply.Code = errorCode.Code
	reply.Msg = errorCode.Msg
	tradelog.Logger.Infof("Receive an CloseTrade request: symbol = %v reply = %v", symbol, reply.String())

	// 数据点的fields和tags
	latency := time.Since(time_start)
	fields := map[string]interface{}{
		"latency": latency.Nanoseconds(), // 时延
	}

	tags := map[string]string{
		"service": "closetrade",                             // 访问的接口
		"code":    strconv.FormatInt(int64(reply.Code), 10), //  状态码
		"msg":     reply.Msg,                                // 状态消息
		"server":  errcode.ServiceName,
	}
	// 一个数据点（一个请求）
	point, err := influx.NewPoint("requests", tags, fields, time.Now())
	if err != nil {
		tradelog.Logger.Fatalf("NewPoint error: %v", err)
		return &reply, err
	}
	select {
	case <-influxdb.StopChan:
		return &reply, nil
	case influxdb.PointChan <- point:
	}

	return &reply, nil
}
