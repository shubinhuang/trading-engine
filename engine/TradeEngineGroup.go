package engine

// 单线程代码，map中存储撮合引擎实例

type TradeEngineGroup struct {
	Engines map[string]*TradeEngine
}

// 处理多标的
func (teg *TradeEngineGroup) ProcessMultiOrder(order *Order) []TradeResult {
	symbol := order.Symbol
	tradeEngine, ok := teg.Engines[symbol]
	if !ok {
		// tradeEngine = &TradeEngine{Symbol: symbol} // TODO 增加tradeEngine的构造函数
		tradeEngine = NewTradeEngine(symbol, order.Price)
		tradeEngine.Symbol = symbol
		teg.Engines[symbol] = tradeEngine
	}
	return tradeEngine.CreateOrder(order)
}

func NewTradeEngineGroup() *TradeEngineGroup {
	engines := make(map[string]*TradeEngine, 0)
	return &TradeEngineGroup{Engines: engines}
}
