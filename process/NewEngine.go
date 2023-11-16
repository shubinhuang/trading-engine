package process

import (
	"trading-engine/engine"
	"trading-engine/errcode"
	"trading-engine/rediscache"

	"github.com/shopspring/decimal"
)

// 开启指定标的的撮合，即在OrderChanMap中增加一个channel，指定标的的order放到该channel中
func NewEngine(symbol string, marketPrice decimal.Decimal) *errcode.ErrCode {
	// if engine.OrderChanMap[symbol] != nil {
	// 	// 该标的的撮合引擎已经打开
	// 	return errcode.EngineExist
	// }

	_, ok := engine.OrderChanMap.Load(symbol)
	if ok {
		// 该标的的撮合引擎已经打开
		return errcode.EngineExist
	}

	// 打开该标的的撮合引擎，
	// 即创建一个channel存放指定标的的order，chan队列长度为500
	// 创建一个协程运行撮合引擎
	engine.OrderChanMap.Store(symbol, make(chan engine.Order, 10000))
	go engine.Run(symbol, marketPrice)

	// 缓存开启的撮合 (symbol, price)
	rediscache.SaveSymbol(symbol)
	rediscache.SavePrice(symbol, marketPrice)

	return errcode.OK
}
