package process

import (
	"trading-engine/engine"
	"trading-engine/errcode"
)

// 关闭指定标的的撮合引擎，即关闭指定标的的order channel
func CloseEngine(symbol string) *errcode.ErrCode {
	// if engine.OrderChanMap[symbol] == nil {
	// 	// 该标的的撮合引擎未开启
	// 	return errcode.EngineNotFound
	// }

	orderChan, ok := engine.OrderChanMap.Load(symbol)
	if ok == false {
		// 该标的的撮合引擎未开启
		return errcode.EngineNotFound
	}

	// 关闭指定标的的order channel
	// 该交易对不会再接收新的委托单，撮合引擎处理完所有委托单后将清除缓存后退出
	// close(engine.OrderChanMap[symbol])
	close(orderChan.(chan engine.Order))
	// engine.OrderChanMap.Delete(symbol)

	return errcode.OK
}
