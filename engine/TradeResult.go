package engine

import (
	"encoding/json"
	"trading-engine/enum"

	"github.com/shopspring/decimal"
)

// 一次交易记录
type TradeResult struct {
	Price          decimal.Decimal     `json:"price"`           // 交易价格
	Quantity       decimal.Decimal     `json:"quantity"`        // 交易数量
	TakerOrderId   string              `json:"taker_orderid"`   // 吃单Id
	MakerOrderId   string              `json:"maker_orderid"`   //	挂单Id
	TakerDirection enum.OrderDirection `json:"taker_direction"` //  吃单的方向（买/卖）
	TimeStamp      int64               `json:"timestamp"`       // 成交时间
}

func (tres *TradeResult) FromJSON(msg []byte) error {
	return json.Unmarshal(msg, tres)
}

func (tres *TradeResult) ToJSON() []byte {
	str, _ := json.Marshal(tres)
	return str
}

func (tres *TradeResult) ToMap() map[string]interface{} {
	data, _ := json.Marshal(tres)
	m := make(map[string]interface{})
	json.Unmarshal(data, &m)
	return m
}

// type TradeResult struct {
// 	TakerOrder   Order     // 吃单的
// 	TradeDetails list.List // 该委托单产生的一系列成交记录
// }

// func (tradeResult *TradeResult) Add(price decimal.Decimal, quantity decimal.Decimal, makerOrder Order, timestamp int64) {
// 	record := TradeRecord{Price: price, Quantity: quantity, TakerOrderId: tradeResult.TakerOrder.OrderId, MakerOrderId: makerOrder.OrderId}
// 	tradeResult.TradeDetails.PushBack(record)
// }
