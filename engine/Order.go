package engine

import (
	"encoding/json"
	"strconv"
	"trading-engine/enum"

	"github.com/shopspring/decimal"
)

// 委托单结构体
type Order struct {

	// 订单ID / 用户ID / 交易对ID:
	OrderId string `json:"order_id"`
	UserId  string `json:"user_id"`
	Symbol  string `json:"symbol"`

	// 动作 / 价格 / 方向 / 状态:
	Action    enum.OrderAction    `json:"action"`
	Price     decimal.Decimal     `json:"price"`
	Direction enum.OrderDirection `json:"direction"`
	Status    enum.OrderStatus    `json:"status"`

	// 订单数量 / 未成交数量:
	Quantity         decimal.Decimal `json:"quantity"`
	UnfilledQuantity decimal.Decimal `json:"unfilled_quantity"`

	// 创建和更新时间:
	CreateAt int64 `json:"create_at"`
	UpdateAt int64 `json:"update_at"`
}

// 更新委托单
// 委托单创建后，由撮合引擎处理时，只有unfilledQuantity、status、updateAt会发生变化，其他属性均为只读。
func (order *Order) UpdateOrder(unfilledQuantity decimal.Decimal, status enum.OrderStatus, updateAt int64) {
	order.UnfilledQuantity = unfilledQuantity
	order.Status = status
	order.UpdateAt = updateAt
}

func (order *Order) FromJSON(msg []byte) error {
	return json.Unmarshal(msg, order)
}

func (order *Order) ToJSON() []byte {
	str, _ := json.Marshal(order)
	return str
}

func (order *Order) ToMap() map[string]interface{} {
	data, _ := json.Marshal(order)
	m := make(map[string]interface{})
	json.Unmarshal(data, &m)
	return m
}

func (order *Order) FromMap(m map[string]string) error {
	var err error
	order.OrderId = m["order_id"]
	order.UserId = m["user_id"]
	order.Symbol = m["symbol"]
	order.Action = enum.OrderAction(m["action"])
	order.Price, err = decimal.NewFromString(m["price"])
	if err != nil {
		// order.Price = decimal.Zero
		return err
	}
	order.Direction = enum.OrderDirection(m["direction"])
	order.Status = enum.OrderStatus(m["status"])
	order.Quantity, err = decimal.NewFromString(m["quantity"])
	if err != nil {
		// order.Quantity = decimal.Zero
		return err
	}
	order.UnfilledQuantity, err = decimal.NewFromString(m["unfilled_quantity"])
	if err != nil {
		// order.UnfilledQuantity = decimal.Zero
		return err
	}
	order.CreateAt, err = strconv.ParseInt(m["create_at"], 10, 64)
	if err != nil {
		// order.CreateAt = 0
		return err
	}
	order.UpdateAt, err = strconv.ParseInt(m["update_at"], 10, 64)
	if err != nil {
		// order.UpdateAt = 0
		return err
	}
	return nil
}
