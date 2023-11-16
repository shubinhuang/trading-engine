package process

import (
	"log"
	"strings"
	"trading-engine/engine"
	"trading-engine/rediscache"
)

// 启动服务时，从缓存中读取程序关闭前开启撮合的标的，进行数据恢复

func Init() {
	symbols := rediscache.GetSymbols()
	log.Printf("%+v\n", symbols)
	for _, symbol := range symbols {
		marketPrice := rediscache.GetPrice(symbol)
		NewEngine(symbol, marketPrice)

		// 数据恢复：读取缓存中委托单，放入对应channel中，重新撮合，恢复到程序关闭前的状态
		orderIdsWithAction := rediscache.GetOrderIdsWithAction(symbol)
		for _, orderId_action := range orderIdsWithAction {
			s := strings.Split(orderId_action, "+")
			orderId := s[0]
			action := s[1]
			orderMap := rediscache.GetOrder(symbol, orderId, action)
			order := engine.Order{}
			err := order.FromMap(orderMap)
			if err != nil {
				log.Print(err)
				continue
			}
			// engine.OrderChanMap[order.Symbol] <- order
			orderChan, _ := engine.OrderChanMap.Load(order.Symbol)
			orderChan.(chan engine.Order) <- order
		}
	}
}
