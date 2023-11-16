package engine

import "sync"

// var OrderChanMap map[string]chan Order
var OrderChanMap sync.Map

func Init() {
	// OrderChanMap = make(map[string]chan Order)
	OrderChanMap = sync.Map{}
}
