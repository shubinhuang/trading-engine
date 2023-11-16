package main

import (
	"fmt"
	"hash/fnv"
	"trading-engine/engine"
	"trading-engine/enum"
	"trading-engine/rediscache"

	"github.com/emirpasic/gods/maps/treemap"
	"github.com/go-redis/redis"
	"github.com/shopspring/decimal"
)

type User struct {
	id   int
	name string
}

// 升序
func byID(a, b interface{}) int {

	// Type assertion, program will panic if this is not respected
	c1 := a.(User)
	c2 := b.(User)

	switch {
	case c1.id > c2.id:
		return 1
	case c1.id < c2.id:
		return -1
	default:
		return 0
	}
}

func testGods() {
	m := treemap.NewWithIntComparator() // empty (keys are of type int)
	m.Put(1, "x")                       // 1->x
	m.Put(2, "b")                       // 1->x, 2->b (in order)
	m.Put(1, "a")                       // 1->a, 2->b (in order)
	k, ok := m.Get(2)                   // b, true
	_, _ = m.Get(3)                     // nil, false
	_ = m.Values()                      // []interface {}{"a", "b"} (in order)
	_ = m.Keys()                        // []interface {}{1, 2} (in order)
	m.Remove(1)                         // 2->b
	m.Clear()                           // empty
	m.Empty()                           // true
	m.Size()                            // 0

	// Other:
	m.Min() // Returns the minimum key and its value from map.
	m.Max() // Returns the maximum key and its value from map.
	fmt.Println(k, ok)

	userMap := treemap.NewWith(byID)
	userMap.Put(User{2, "Se"}, "x")
	userMap.Put(User{3, "Sec"}, "y")
	userMap.Put(User{1, "S"}, 5.6)
	fmt.Println(userMap.Keys()...)
	fmt.Println(userMap.Values()...)
	a, b := userMap.Min()
	fmt.Println(a, b)
	userMap.Remove(User{5, "s"})

}

func testEngine() {
	ob := engine.OrderBook{Direction: enum.SELL, Book: treemap.NewWith(engine.SortSell)}
	ob.Add(engine.Order{Price: decimal.NewFromInt(123)})
	ob.Add(engine.Order{Price: decimal.NewFromInt(321)})
	fmt.Printf("%+v\n", ob.GetFirst())
	ob.Remove(&engine.Order{Price: decimal.NewFromInt(123)})
	fmt.Printf("%+v\n", ob.GetFirst())
}

func testRedis() {
	rediscache.RedisInit()
	rdc := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	fmt.Println(rdc.Exists("aaaa").Val())
}

// func test() {

// 	// testRedis()
// 	// testGods()

// 	// 测试撮合引擎
// 	testSingle()
// 	testMulti()
// 	testCancelOrder()

// 	orderIdsWithAction := []string{
// 		"123" + "+" + "sss",
// 		"456" + "+" + "ttt",
// 		"789" + "+" + "hhh",
// 	}

// 	for _, orderId_action := range orderIdsWithAction {
// 		s := strings.Split(orderId_action, "+")
// 		orderId := s[0]
// 		action := s[1]
// 		fmt.Println(orderId, action)

// 	}
// }

func hash(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}

// func main() {
// 	fmt.Println(hash("s1") % 32)
// 	fmt.Println(hash("s5") % 32)
// }
