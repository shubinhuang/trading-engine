package engine

import (
	"log"
	"time"
	"trading-engine/enum"
	"trading-engine/errcode"
	"trading-engine/influxdb"
	"trading-engine/rediscache"

	influx "github.com/influxdata/influxdb1-client/v2"
	"github.com/shopspring/decimal"
)

type TradeEngine struct {
	Symbol      string          // 哪个标的的撮合引擎
	BuyBook     OrderBook       //买单队列
	SellBook    OrderBook       // 卖单队列
	MarketPrice decimal.Decimal // 最新市场价
}

func NewTradeEngine(symbol string, marketPrice decimal.Decimal) *TradeEngine {
	return &TradeEngine{
		Symbol:      symbol,
		BuyBook:     NewOrderBook(enum.BUY),
		SellBook:    NewOrderBook(enum.SELL),
		MarketPrice: marketPrice,
	}
}

// 撤单，action==cancel
func (tradeEngine *TradeEngine) CancelOrder(order *Order) bool {
	var cancelSuccess bool
	// 从orderBook中删除该订单
	switch order.Direction {
	case enum.BUY:
		cancelSuccess = tradeEngine.BuyBook.DeleteFromBook(order)
	case enum.SELL:
		cancelSuccess = tradeEngine.SellBook.DeleteFromBook(order)
	}
	// 从缓存中删除该订单
	rediscache.RemoveOrder(order.Symbol, order.OrderId, enum.CREATE_ORDER.String())
	// 同时删除这个撤单请求的缓存
	rediscache.RemoveOrder(order.Symbol, order.OrderId, enum.CANCEL_ORDER.String())

	// 返回撤单结果
	return cancelSuccess
}

// 下单，action==create
func (tradeEngine *TradeEngine) CreateOrder(order *Order) []TradeResult {
	switch order.Direction {
	case enum.BUY: // 买单  与SellBook中的订单撮合，剩余的放入BuyBook
		// fmt.Println("Process buy")
		return tradeEngine.CreateAnOrder(order, &tradeEngine.SellBook, &tradeEngine.BuyBook)
	case enum.SELL: // 卖单 与BuyBook中的订单撮合，剩余的放入SellBook
		// fmt.Println("Process sell")
		return tradeEngine.CreateAnOrder(order, &tradeEngine.BuyBook, &tradeEngine.SellBook)
	default:
		return nil
		// 返回撮合结果
	}
}

// 处理买单&处理卖单
func (tradeEngine *TradeEngine) CreateAnOrder(takerOrder *Order, makerBook *OrderBook, anotherBook *OrderBook) []TradeResult {
	timestamp := takerOrder.CreateAt
	tradeResults := make([]TradeResult, 0)

	takerUnfilledQuantity := takerOrder.Quantity
	if !takerOrder.UnfilledQuantity.Equal(decimal.Zero) {
		takerUnfilledQuantity = takerOrder.UnfilledQuantity
	}

	for {
		makerOrderP := makerBook.GetFirst() // 对手盘排最前面的一单
		if makerOrderP == nil {
			break
		}
		makerOrder := makerOrderP.(Order)
		if takerOrder.Direction == enum.BUY && takerOrder.Price.LessThan(makerOrder.Price) {
			break
		} else if takerOrder.Direction == enum.SELL && takerOrder.Price.GreaterThan(makerOrder.Price) {
			break
		}

		// 发生成交
		var tradeQuantity decimal.Decimal

		// 交易量
		if takerUnfilledQuantity.LessThanOrEqual(makerOrder.UnfilledQuantity) {
			tradeQuantity = takerUnfilledQuantity
		} else {
			tradeQuantity = makerOrder.UnfilledQuantity
		}

		tradeEngine.MarketPrice = makerOrder.Price
		// 更新缓存中交易对的最新成交价
		rediscache.SavePrice(tradeEngine.Symbol, tradeEngine.MarketPrice)

		// 一条成交记录
		tradeResults = append(tradeResults, TradeResult{
			Price:          makerOrder.Price,
			Quantity:       tradeQuantity,
			TakerOrderId:   takerOrder.OrderId,
			MakerOrderId:   makerOrder.OrderId,
			TakerDirection: takerOrder.Direction,
			TimeStamp:      timestamp,
		})

		// 往influxDB中记录某某时刻成交一笔
		// 数据点的fields和tags
		tradePrice, _ := makerOrder.Price.Float64()
		fields := map[string]interface{}{
			"price": tradePrice, // 成交价
		}
		tags := map[string]string{
			"symbol": tradeEngine.Symbol, // 标的
			"server": errcode.ServiceName,
		}
		// 一个数据点（一个请求）
		point, err := influx.NewPoint("trade_record", tags, fields, time.Now())
		if err != nil {
			log.Fatalf("NewPoint error: %v", err)
		}

		select {
		case <-influxdb.StopChan:
		case influxdb.PointChan <- point:
		}

		takerUnfilledQuantity = takerUnfilledQuantity.Sub(tradeQuantity)

		makerUnfilledQuantity := makerOrder.UnfilledQuantity.Sub(tradeQuantity)

		if makerUnfilledQuantity.IsZero() {
			// 挂单完全成交，从orderBook中删除
			makerOrder.UpdateOrder(makerUnfilledQuantity, enum.FULLY_FILLED, timestamp)
			makerBook.Add(makerOrder)
			makerBook.Remove(&makerOrder)
			// 从缓存中删除该订单
			rediscache.RemoveOrder(makerOrder.Symbol, makerOrder.OrderId, enum.CREATE_ORDER.String())
		} else {
			// 挂单部分成交，更新orderBook中的挂单，更新缓存中的挂单
			makerOrder.UpdateOrder(makerUnfilledQuantity, enum.PARTIAL_FILLED, timestamp)
			makerBook.Add(makerOrder)
			rediscache.UpdateOrder(makerOrder.ToMap())
		}

		if takerUnfilledQuantity.IsZero() {
			// 吃单完全成交，不放入orderBook
			takerOrder.UpdateOrder(takerUnfilledQuantity, enum.FULLY_FILLED, timestamp)
			// （不更新缓存）直接删除缓存中的吃单
			// rediscache.UpdateOrder(takerOrder.ToMap())
			rediscache.RemoveOrder(takerOrder.Symbol, takerOrder.OrderId, enum.CREATE_ORDER.String())
			break
		}

	}

	if takerUnfilledQuantity.GreaterThan(decimal.Zero) {
		// 吃单部分成交或者未成交，更新吃单，更新缓存中的吃单，将吃单加入orderBook

		if takerUnfilledQuantity == takerOrder.Quantity {
			takerOrder.UpdateOrder(takerUnfilledQuantity, enum.PENDING, timestamp)
		} else {
			takerOrder.UpdateOrder(takerUnfilledQuantity, enum.PARTIAL_FILLED, timestamp)
		}

		rediscache.UpdateOrder(takerOrder.ToMap())
		anotherBook.Add(*takerOrder)
	}

	return tradeResults
}
