package main

import (
	"log"
	"strconv"
	"trading-engine/consul"
	"trading-engine/engine"
	"trading-engine/process"
	"trading-engine/rediscache"
)

func init() {
	engine.Init()
	rediscache.RedisInit()
	process.Init()
}

func main() {
	// testSingle()
	services := consul.DiscoveryService("trade-0", "")
	var addr string
	if services[0].Service.Address == "" {
		addr = services[0].Node.Address + ":" + strconv.Itoa(services[0].Service.Port)
	} else {
		addr = services[0].Service.Address + ":" + strconv.Itoa(services[0].Service.Port)
	}
	log.Printf("%+v", addr)
}
