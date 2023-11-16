package rediscache

import (
	"log"
	"strings"
	"trading-engine/errcode"

	"github.com/go-redis/redis"
	"github.com/shopspring/decimal"
)

// 缓存开启撮合的交易对
func SaveSymbol(symbol string) {
	key := errcode.ServiceName + ":symbols"
	// 保存到set类型
	RedisClient.SAdd(key, symbol)
}

// 关闭撮合时移除相应交易对缓存
func RemoveSymbol(symbol string) {
	key := errcode.ServiceName + ":symbols"
	RedisClient.SRem(key, symbol)
}

func GetSymbols() []string {
	key := errcode.ServiceName + ":symbols"
	log.Printf(key)
	return RedisClient.SMembers(key).Val()
}

// 缓存交易对的最新市场价
func SavePrice(symbol string, MarketPrice decimal.Decimal) {
	key := errcode.ServiceName + ":price:" + symbol
	RedisClient.Set(key, MarketPrice.String(), 0)
}

func RemovePrice(symbol string) {
	key := errcode.ServiceName + ":price:" + symbol
	RedisClient.Del(key)
}

func GetPrice(symbol string) decimal.Decimal {
	key := errcode.ServiceName + ":price:" + symbol
	priceStr := RedisClient.Get(key).Val()
	result, err := decimal.NewFromString(priceStr)
	if err != nil {
		result = decimal.Zero
	}
	return result
}

// 缓存委托单
func SaveOrder(order map[string]interface{}) {
	// order.ToMap()
	symbol := order["symbol"].(string)
	orderId := order["order_id"].(string)
	timestamp := order["create_at"].(float64)
	action := order["action"].(string)

	// 缓存一个完整订单
	key := errcode.ServiceName + ":order:" + symbol + ":" + orderId + ":" + action
	RedisClient.HMSet(key, order)

	// 缓存一个交易对的所有订单，按请求顺序缓存, sorted set
	key = errcode.ServiceName + ":orderids:" + symbol
	z := redis.Z{
		Score:  timestamp,
		Member: orderId + "+" + action,
	}
	RedisClient.ZAdd(key, z)
}

// 获得该交易对下所有的委托单的orderId及其action  orderId + "+" + action
func GetOrderIdsWithAction(symbol string) []string {
	key := errcode.ServiceName + ":orderids:" + symbol
	return RedisClient.ZRange(key, 0, -1).Val()
}

// 缓存中是否存在该委托单，用于去重
func OrderExist(symbol, orderId, action string) bool {
	key := errcode.ServiceName + ":order:" + symbol + ":" + orderId + ":" + action
	if RedisClient.Exists(key).Val() > 0 {
		return true
	} else {
		return false
	}
}

func GetOrder(symbol, orderId, action string) map[string]string {
	key := errcode.ServiceName + ":order:" + symbol + ":" + orderId + ":" + action
	orderMap := RedisClient.HGetAll(key).Val()
	// tmpO := engine.Order{}
	// orderMap := tmpO.ToMap()
	// for field := range orderMap {
	// 	orderMap[key] = RedisClient.HGet(key, field).Val()
	// }

	return orderMap
}

// 删除委托单
func RemoveOrder(symbol string, orderId string, action string) {
	// symbol := orderMap["symbol"].(string)
	// orderId := orderMap["order_id"].(string)
	// // timestamp := orderMap["create_at"].(float64)
	// action := orderMap["action"].(string)
	key := errcode.ServiceName + ":order:" + symbol + ":" + orderId + ":" + action
	RedisClient.Del(key)

	key = errcode.ServiceName + ":orderids:" + symbol
	RedisClient.ZRem(key, orderId+"+"+action)

}

// 修改缓存中的订单
func UpdateOrder(orderMap map[string]interface{}) {
	symbol := orderMap["symbol"].(string)
	orderId := orderMap["order_id"].(string)
	action := orderMap["action"].(string)
	key := errcode.ServiceName + ":order:" + symbol + ":" + orderId + ":" + action
	RedisClient.HMSet(key, orderMap)
}

// 删除该交易对下所有委托单
func ClearSymbol(symbol string) {
	orderIdsWithAction := GetOrderIdsWithAction(symbol)

	for _, orderId_action := range orderIdsWithAction {
		s := strings.Split(orderId_action, "+")
		orderId := s[0]
		action := s[1]
		// 删除缓存中的一个委托单对象
		key := errcode.ServiceName + ":order:" + symbol + ":" + orderId + ":" + action
		RedisClient.Del(key)
	}
	// 删除缓存中的该交易对下的orderId集合
	key := errcode.ServiceName + ":orderids:" + symbol
	RedisClient.Del(key)

	// 清除交易对及其市场价的缓存
	RemoveSymbol(symbol)
	RemovePrice(symbol)
}

// 往Redis的一个stream添加一个撤单结果
func CancelResult(symbol string, orderId string, isSuccess bool) {
	values := map[string]interface{}{"orderId": orderId, "isSuccess": isSuccess}

	x := &redis.XAddArgs{
		Stream:       errcode.ServiceName + ":cancelresults:" + symbol,
		MaxLenApprox: 1000,
		Values:       values,
	}
	RedisClient.XAdd(x)
}

// 往Redis的一个stream添加一个撮合结果
func TradeResult(symbol string, tradeResult map[string]interface{}) {
	x := &redis.XAddArgs{
		Stream:       errcode.ServiceName + ":traderesults:" + symbol,
		MaxLenApprox: 10000,
		Values:       tradeResult,
	}
	RedisClient.XAdd(x)
}
