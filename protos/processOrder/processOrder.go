package processOrder

import (
	context "context"
	"strconv"
	"time"
	"trading-engine/engine"
	"trading-engine/enum"
	"trading-engine/errcode"
	"trading-engine/influxdb"
	"trading-engine/process"
	"trading-engine/tradelog"

	influx "github.com/influxdata/influxdb1-client/v2"
	"github.com/shopspring/decimal"
)

type Server struct {
	UnimplementedOrderServiceServer
}

func (server *Server) CreateOrder(ctx context.Context, in *CreateRequest) (*CreateReply, error) {
	time_start := time.Now()
	reply := CreateReply{}
	symbol := in.GetSymbol()
	userId := in.GetUserId()
	orderId := in.GetOrderId()
	direction := in.GetDirection()
	price, err := decimal.NewFromString(in.GetPrice())

	if err != nil {
		reply.Code = errcode.FORMATERROR.Code
		reply.Msg = errcode.FORMATERROR.Msg

		return &reply, err
	}

	quantity, err := decimal.NewFromString(in.GetQuantity())

	if err != nil {
		reply.Code = errcode.FORMATERROR.Code
		reply.Msg = errcode.FORMATERROR.Msg
		return &reply, err
	}

	order := engine.Order{
		Symbol:    symbol,
		UserId:    userId,
		OrderId:   orderId,
		Action:    enum.CREATE_ORDER,
		Direction: enum.OrderDirection(direction),
		Price:     price,
		Quantity:  quantity,
	}

	errorCode := process.Dispatch(order)
	reply.Code = errorCode.Code
	reply.Msg = errorCode.Msg
	// reply := CreateReply{Code: errorCode.Code, Msg: errorCode.Msg}
	tradelog.Logger.Infof("Receive an CreateOrder request: symbol = %v orderId = %v reply = %v", symbol, orderId, reply.String())

	// 数据点的fields和tags
	latency := time.Since(time_start)
	fields := map[string]interface{}{
		"latency": latency.Nanoseconds(), // 时延
	}

	tags := map[string]string{
		"service": "createorder",                            // 访问的接口
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

func (server *Server) CancelOrder(ctx context.Context, in *CancelRequest) (*CancelReply, error) {
	time_start := time.Now()

	symbol := in.GetSymbol()
	orderId := in.GetOrderId()
	direction := in.GetDirection()
	order := engine.Order{
		Symbol:    symbol,
		OrderId:   orderId,
		Action:    enum.CANCEL_ORDER,
		Direction: enum.OrderDirection(direction),
	}
	errorCode := process.Dispatch(order)
	reply := CancelReply{Code: errorCode.Code, Msg: errorCode.Msg}
	tradelog.Logger.Infof("Receive an CancelOrder request: symbol = %v orderId = %v reply = %v", symbol, orderId, reply.String())

	// 数据点的fields和tags
	latency := time.Since(time_start)
	fields := map[string]interface{}{
		"latency": latency.Nanoseconds(), // 时延
	}

	tags := map[string]string{
		"service": "cancelorder",                                // 访问的接口
		"code":    strconv.FormatInt(int64(errorCode.Code), 10), //  状态码
		"msg":     errorCode.Msg,                                // 状态消息
		"server":  errcode.ServiceName,
	}
	// 一个数据点（一个请求）
	point, err := influx.NewPoint("requests", tags, fields, time.Now())
	if err != nil {
		tradelog.Logger.Fatalf("NewPoint error: %v", err)
	}

	select {
	case <-influxdb.StopChan:
		return &reply, nil
	case influxdb.PointChan <- point:
	}

	return &reply, nil
}
