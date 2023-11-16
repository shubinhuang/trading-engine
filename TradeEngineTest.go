package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
	engine "trading-engine/engine"
	"trading-engine/enum"

	"github.com/shopspring/decimal"
)

func CreateOrder(symbol string, timestamp int64, orderId string, userId string, direction enum.OrderDirection, price decimal.Decimal, quantity decimal.Decimal) engine.Order {

	order := engine.Order{}
	order.OrderId = orderId
	order.UserId = userId
	order.Symbol = symbol
	order.Direction = direction
	order.Price = price
	order.Quantity = quantity
	order.UnfilledQuantity = quantity
	order.CreateAt = timestamp
	order.UpdateAt = timestamp

	return order
}

type testOrderInfo struct {
	symbol    string
	direction enum.OrderDirection
	price     decimal.Decimal
	quantity  decimal.Decimal
}

func testSingle() {
	fmt.Println("单标的撮合模拟")
	input := []testOrderInfo{
		{"0", enum.BUY, decimal.NewFromFloat(2082.34), decimal.NewFromInt(1)},
		{"0", enum.SELL, decimal.NewFromFloat(2087.60), decimal.NewFromInt(2)},
		{"0", enum.BUY, decimal.NewFromFloat(2087.80), decimal.NewFromInt(1)},
		{"0", enum.BUY, decimal.NewFromFloat(2085.01), decimal.NewFromInt(5)},
		{"0", enum.SELL, decimal.NewFromFloat(2088.02), decimal.NewFromInt(3)},
		{"0", enum.SELL, decimal.NewFromFloat(2087.60), decimal.NewFromInt(6)},
		{"0", enum.BUY, decimal.NewFromFloat(2081.11), decimal.NewFromInt(7)},
		{"0", enum.BUY, decimal.NewFromFloat(2086.00), decimal.NewFromInt(3)},
		{"0", enum.BUY, decimal.NewFromFloat(2088.33), decimal.NewFromInt(1)},
		{"0", enum.SELL, decimal.NewFromFloat(2086.54), decimal.NewFromInt(2)},
		{"0", enum.SELL, decimal.NewFromFloat(2086.55), decimal.NewFromInt(5)},
		{"0", enum.BUY, decimal.NewFromFloat(2086.55), decimal.NewFromInt(3)},
	}

	tradeEngine := engine.NewTradeEngine("0", decimal.Zero)

	for _, o := range input {
		timestamp := time.Now().UnixMicro()
		userId := strconv.Itoa(rand.Int())
		orderId := strconv.Itoa(rand.Int())
		order := CreateOrder(o.symbol, timestamp, orderId, userId, o.direction, o.price, o.quantity)

		tradeResult := tradeEngine.CreateOrder(&order)

		// fmt.Println("Trade result")
		if tradeResult != nil {
			for _, res := range tradeResult {
				fmt.Printf("%+v\n", res)
			}
		}
	}

	fmt.Println("\n\nAll deal down")
	fmt.Println("Buy order queue")
	for _, bo := range tradeEngine.BuyBook.Book.Values() {
		buyOrder := bo.(engine.Order)
		fmt.Printf("%v\t%v\n", buyOrder.Price, buyOrder.UnfilledQuantity)
	}
	fmt.Println("------------")
	fmt.Println("Sell order queue")
	for _, so := range tradeEngine.SellBook.Book.Values() {
		sellOrder := so.(engine.Order)
		fmt.Printf("%v\t%v\n", sellOrder.Price, sellOrder.UnfilledQuantity)
	}
	fmt.Println("------------")
	fmt.Printf("Final market price \n%v\n", tradeEngine.MarketPrice)
}

func testMulti() {
	fmt.Println("多标的撮合模拟")
	input := []testOrderInfo{
		{"1", enum.BUY, decimal.NewFromFloat(2082.34), decimal.NewFromInt(1)},
		{"1", enum.SELL, decimal.NewFromFloat(2087.60), decimal.NewFromInt(2)},
		{"1", enum.BUY, decimal.NewFromFloat(2087.80), decimal.NewFromInt(1)},
		{"1", enum.BUY, decimal.NewFromFloat(2085.01), decimal.NewFromInt(5)},
		{"1", enum.SELL, decimal.NewFromFloat(2088.02), decimal.NewFromInt(3)},
		{"1", enum.SELL, decimal.NewFromFloat(2087.60), decimal.NewFromInt(6)},
		{"2", enum.BUY, decimal.NewFromFloat(2081.11), decimal.NewFromInt(7)},
		{"2", enum.BUY, decimal.NewFromFloat(2086.00), decimal.NewFromInt(3)},
		{"2", enum.BUY, decimal.NewFromFloat(2088.33), decimal.NewFromInt(1)},
		{"2", enum.SELL, decimal.NewFromFloat(2086.54), decimal.NewFromInt(2)},
		{"2", enum.SELL, decimal.NewFromFloat(2086.55), decimal.NewFromInt(5)},
		{"2", enum.BUY, decimal.NewFromFloat(2086.55), decimal.NewFromInt(3)},
	}

	tradeEngineGroup := engine.NewTradeEngineGroup()

	for _, o := range input {
		timestamp := time.Now().UnixMicro()
		userId := strconv.Itoa(rand.Int())
		orderId := strconv.Itoa(rand.Int())
		order := CreateOrder(o.symbol, timestamp, orderId, userId, o.direction, o.price, o.quantity)
		tradeResult := tradeEngineGroup.ProcessMultiOrder(&order)

		// fmt.Println("Trade result")
		if tradeResult != nil {
			for _, res := range tradeResult {
				fmt.Printf("%+v\n", res)
			}
		}
	}

}

func testCancelOrder() {
	fmt.Println("测试撤单")
	input := []testOrderInfo{
		{"1", enum.BUY, decimal.NewFromFloat(2082.34), decimal.NewFromInt(1)},
		{"1", enum.BUY, decimal.NewFromFloat(2087.60), decimal.NewFromInt(2)},
		{"1", enum.BUY, decimal.NewFromFloat(2087.80), decimal.NewFromInt(1)},
		{"1", enum.BUY, decimal.NewFromFloat(2085.01), decimal.NewFromInt(5)},
		{"1", enum.BUY, decimal.NewFromFloat(2088.02), decimal.NewFromInt(3)},
		{"1", enum.BUY, decimal.NewFromFloat(2087.60), decimal.NewFromInt(6)},
	}

	tradeEngine := engine.NewTradeEngine("1", decimal.Zero)

	orders := make([]engine.Order, 0)
	for _, o := range input {
		timestamp := time.Now().UnixMicro()
		userId := strconv.Itoa(rand.Int())
		orderId := strconv.Itoa(rand.Int())
		order := CreateOrder(o.symbol, timestamp, orderId, userId, o.direction, o.price, o.quantity)
		orders = append(orders, order)
		tradeEngine.CreateOrder(&order)
	}

	fmt.Println("Buy order queue before cancel order")
	for _, bo := range tradeEngine.BuyBook.Book.Values() {
		buyOrder := bo.(engine.Order)
		fmt.Printf("%v\t%v\n", buyOrder.Price, buyOrder.UnfilledQuantity)
	}
	fmt.Println("------------")

	tradeEngine.CancelOrder(&orders[3])

	fmt.Println("Buy order queue after cancel order")
	for _, bo := range tradeEngine.BuyBook.Book.Values() {
		buyOrder := bo.(engine.Order)
		fmt.Printf("%v\t%v\n", buyOrder.Price, buyOrder.UnfilledQuantity)
	}
	fmt.Println("------------")

}

// func main() {
// 	testSingle()
// 	// testMulti()
// 	// testCancelOrder()
// }
