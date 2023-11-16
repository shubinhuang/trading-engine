package rediscache

import (
	"fmt"
	"strconv"
	"trading-engine/consul"

	"github.com/go-redis/redis"
)

var RedisClient *redis.Client

func RedisInit() {
	services := consul.DiscoveryService("redis", "")

	var addr = services[0].Node.Address + ":" + strconv.Itoa(services[0].Service.Port)

	// log.Printf("%+v", addr)
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "",
		DB:       0,
	})

	res, err := RedisClient.Ping().Result()
	if err != nil {
		panic(err)
	} else {
		fmt.Println(res)
	}
}
