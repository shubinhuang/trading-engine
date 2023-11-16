package consul

import (
	"log"

	consual_api "github.com/hashicorp/consul/api"
)

var consulClient *consual_api.Client

func InitConsul() {
	consulConfig := consual_api.DefaultConfig()
	consulConfig.Address = "192.168.78.128:8500"
	var err error
	consulClient, err = consual_api.NewClient(consulConfig)
	if err != nil {
		log.Fatalf("new consual client is failed, err: %v", err)
	}
}

// 服务发现
func DiscoveryService(service, tag string) []*consual_api.ServiceEntry {
	services, _, err := consulClient.Health().Service(service, tag, true, nil)

	if err != nil {
		log.Fatalf("get services is failed, err: %v", err)
	}
	return services
}
