package engine

import (
	"strconv"
	"trading-engine/enum"
	"trading-engine/rediscache"

	"github.com/emirpasic/gods/maps/treemap"
	"github.com/emirpasic/gods/utils"
	"github.com/shopspring/decimal"
)

// 价格和时间戳可能相同，orderId是全局唯一的，引入orderId作为key避免key重复，同时便于撤单
type OrderKey struct {
	CreateAt int64
	Price    decimal.Decimal
	OrderId  string
}

// comparator返回1表示ab逆序，需要调换位置，返回-1表示ab顺序，返回0表示key相同
// 卖单队列 按价格递增  价格越低排越前面
func SortSell(a, b interface{}) int {
	c1 := a.(OrderKey)
	c2 := b.(OrderKey)

	if c1.Price.GreaterThan(c2.Price) {
		// 卖单队列，价格低排前面
		return 1
	} else if c1.Price.LessThan(c2.Price) {
		return -1
	} else {
		// 价格相等时，时间早排前面
		if c1.CreateAt > c2.CreateAt {
			return 1
		} else if c1.CreateAt < c2.CreateAt {
			return -1
		}
	}
	// price和时间戳都相等时，按照orderid排序
	return utils.StringComparator(c1.OrderId, c2.OrderId)
}

// 买单队列 按价格递减  价格越高排越前面
func SortBuy(a, b interface{}) int {

	c1 := a.(OrderKey)
	c2 := b.(OrderKey)

	if c1.Price.LessThan(c2.Price) {
		return 1
	} else if c1.Price.GreaterThan(c2.Price) {
		return -1
	} else {
		if c1.CreateAt > c2.CreateAt {
			return 1
		} else if c1.CreateAt < c2.CreateAt {
			return -1
		}
	}
	// price和时间戳都相等时，按照orderid排序
	return utils.StringComparator(c1.OrderId, c2.OrderId)
}

type OrderBook struct {
	// 买单队列or卖单队列
	Direction enum.OrderDirection
	// 用红黑树方式保存
	Book *treemap.Map
}

func NewOrderBook(direction enum.OrderDirection) OrderBook {
	ob := OrderBook{Direction: direction}
	if direction == enum.BUY {
		ob.Book = treemap.NewWith(SortBuy)
	} else if direction == enum.SELL {
		ob.Book = treemap.NewWith(SortSell)
	}
	return ob
}

func (ob *OrderBook) GetFirst() interface{} {
	_, v := ob.Book.Min()
	return v
}

// 根据key从红黑树中删除order
func (ob *OrderBook) Remove(o *Order) {
	ob.Book.Remove(OrderKey{o.CreateAt, o.Price, o.OrderId})
}

// 红黑树中的key是(createAt, price, orderId)，value是order。撤单时直接根据key删除，无需遍历orderBook
func (ob *OrderBook) DeleteFromBook(o *Order) bool {
	// 从缓存中读取欲删除的订单的时间戳和价格
	wanted_order := rediscache.GetOrder(o.Symbol, o.OrderId, enum.CREATE_ORDER.String())

	create_at, err := strconv.ParseInt(wanted_order["create_at"], 10, 64)
	if err != nil {
		return false
	}
	price, err := decimal.NewFromString(wanted_order["price"])
	if err != nil {
		return false
	}

	k := OrderKey{create_at, price, o.OrderId}
	// 欲删除订单存在
	_, ok := ob.Book.Get(k)
	// fmt.Printf("\n%v\n", ob.Book.Keys())
	// fmt.Printf("\n%v  %v\n", k, ok) // 构造的撤单请求没有时间戳和price，导致key错误
	if ok {
		ob.Book.Remove(k)
	}
	return ok
}

// // 红黑树中的key是(createAt, price)，但是撤单时要根据orderId将对应order从orderBook中删除
// // orderId在整个数据库中是唯一的
// func (ob *OrderBook) DeleteFromBook(o *Order) bool {
// 	var ok bool = false
// 	it := ob.Book.Iterator()
// 	it.Begin()
// 	for it.Next() {
// 		// fmt.Printf("\n%v\n", it.Value().(Order))
// 		itOrder := it.Value().(Order)
// 		if o.OrderId == itOrder.OrderId {
// 			ob.Book.Remove(it.Key())
// 			ok = true
// 			break
// 		}
// 	}
// 	return ok
// }

// 往orderBook增加新委托单
func (ob *OrderBook) Add(o Order) {
	ob.Book.Put(OrderKey{o.CreateAt, o.Price, o.OrderId}, o)
}
