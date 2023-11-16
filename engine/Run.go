package engine

import (
	"time"
	"trading-engine/enum"
	"trading-engine/errcode"
	"trading-engine/influxdb"
	"trading-engine/rediscache"

	influx "github.com/influxdata/influxdb1-client/v2"
	"github.com/shopspring/decimal"
)

// 多协程并发，每个协程是一个撮合引擎，从对应的chan中读取order进行处理
func Run(symbol string, marketPrice decimal.Decimal) {
	tradeEngine := NewTradeEngine(symbol, marketPrice)
	orderChan, _ := OrderChanMap.Load(symbol)
	for {
		order, ok := <-orderChan.(chan Order)

		// 写队列时延
		diff := time.Now().UnixMicro() - order.CreateAt
		fields := map[string]interface{}{
			"latency": diff,                        // 队列时延
			"length":  len(orderChan.(chan Order)), //队列长度
		}

		tags := map[string]string{
			"symbol": order.Symbol,
			"server": errcode.ServiceName,
		}
		// 一个数据点（一个请求）
		point, _ := influx.NewPoint("orderchan", tags, fields, time.Now())
		influxdb.PointChan <- point

		// order, ok := <-OrderChanMap[symbol]
		if !ok {
			// 对应的chan已关闭，应关闭该撮合引擎，清空该交易对相应缓存
			// delete(OrderChanMap, symbol)
			OrderChanMap.Delete(symbol)
			rediscache.ClearSymbol(symbol)
			return
		}

		switch order.Action {
		case enum.CREATE_ORDER:
			tradeResults := tradeEngine.CreateOrder(&order)
			if tradeResults != nil {
				for _, res := range tradeResults {
					rediscache.TradeResult(order.Symbol, res.ToMap())
				}
			}

		case enum.CANCEL_ORDER:
			isSuccess := tradeEngine.CancelOrder(&order)
			// 发送撤单结果到redis的一个stream
			rediscache.CancelResult(order.Symbol, order.OrderId, isSuccess)
			// fmt.Printf("symbol: %v orderId: %v cancel result: %v\n", order.Symbol, order.OrderId, isSuccess)
		}
	}

}
