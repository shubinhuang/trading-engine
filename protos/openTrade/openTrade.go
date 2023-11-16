package openTrade

import (
	"context"
	"strconv"
	"time"
	"trading-engine/errcode"
	"trading-engine/influxdb"
	"trading-engine/process"
	"trading-engine/tradelog"

	influx "github.com/influxdata/influxdb1-client/v2"

	"github.com/shopspring/decimal"
)

// 实现服务中定义的OpenTrade方法
type Server struct {
	UnimplementedOpenServiceServer
}

func (sever *Server) OpenTrade(ctx context.Context, in *OpenRequest) (*OpenReply, error) {
	time_start := time.Now()
	reply := OpenReply{}
	symbol := in.GetSymbol()
	marketPrice, err := decimal.NewFromString(in.GetPrice())
	if err != nil {
		reply.Code = errcode.FORMATERROR.Code
		reply.Msg = errcode.FORMATERROR.Msg
		return &reply, err
	}
	errorCode := process.NewEngine(symbol, marketPrice)
	reply.Code = errorCode.Code
	reply.Msg = errorCode.Msg
	// reply := OpenReply{Code: errorCode.Code, Msg: errorCode.Msg}
	tradelog.Logger.Infof("Receive an OpenTrade request: symbol = %v price = %v reply = %v", symbol, marketPrice, reply.String())

	// 数据点的fields和tags
	latency := time.Since(time_start)
	fields := map[string]interface{}{
		"latency": latency.Nanoseconds(), // 时延
	}

	tags := map[string]string{
		"service": "opentrade",                              // 访问的接口
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
