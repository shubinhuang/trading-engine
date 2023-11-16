package process

import (
	"time"
	"trading-engine/engine"
	"trading-engine/enum"
	"trading-engine/errcode"
	"trading-engine/rediscache"
)

// 分发订单到不同的order channel
func Dispatch(order engine.Order) *errcode.ErrCode {
	// if engine.OrderChanMap[order.Symbol] == nil {
	// 	// 该order对应的标的的撮合引擎还未开启
	// 	return errcode.EngineNotFound
	// }

	orderChan, ok := engine.OrderChanMap.Load(order.Symbol)
	if !ok {
		// 该order对应的标的的撮合引擎还未开启
		return errcode.EngineNotFound
	}

	// 判断Order是否存在
	if order.Action == enum.CREATE_ORDER {
		// 重复订单
		if rediscache.OrderExist(order.Symbol, order.OrderId, enum.CREATE_ORDER.String()) {
			return errcode.OrderExist
		}
	} else {
		// 要撤销的委托单不存在
		if !rediscache.OrderExist(order.Symbol, order.OrderId, enum.CREATE_ORDER.String()) {
			return errcode.OrderNotFound
		}
	}

	order.CreateAt = time.Now().UnixMicro()
	order.UpdateAt = order.CreateAt

	// 分发订单到对应的orderChan
	// 如果队列满则要丢弃该请求

	c := orderChan.(chan engine.Order)
	// c <- order
	// // 缓存订单
	// rediscache.SaveOrder(order.ToMap())

	if len(c) > 1000 {
		return errcode.ChanFull
	} else {
		c <- order
		// 缓存订单
		rediscache.SaveOrder(order.ToMap())
	}
	// select {
	// case orderChan.(chan engine.Order) <- order:
	// 	// 缓存订单
	// 	rediscache.SaveOrder(order.ToMap())
	// default:
	// 	return errcode.ChanFull
	// }

	return errcode.OK
}
