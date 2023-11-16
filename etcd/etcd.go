// grpc-etcd/grpc-server/etcd.go

package etcd

import (
	"context"
	"fmt"
	"time"

	eclient "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/naming/endpoints"
)

func RegisterEndPointToEtcd(ctx context.Context, etcdAddr, serverAddr, serverName string) {
	// 创建 etcd 客户端
	etcdClient, _ := eclient.NewFromURL(etcdAddr)
	etcdManager, _ := endpoints.NewManager(etcdClient, serverName)

	// 创建一个租约，每隔 10s 需要向 etcd 汇报一次心跳，证明当前节点仍然存活
	var ttl int64 = 10
	lease, _ := etcdClient.Grant(ctx, ttl)

	// 添加注册节点到 etcd 中，并且携带上租约 id
	// 以 serverName/serverAddr 为 key，serverAddr 为 value
	// serverName/serverAddr 中的 serverAddr 可以自定义，只要能够区分同一个 grpc 服务器功能的不同机器即可

	_ = etcdManager.AddEndpoint(ctx, fmt.Sprintf("%s/%s", serverName, serverAddr), endpoints.Endpoint{Addr: serverAddr}, eclient.WithLease(lease.ID))

	// 每隔 5 s进行一次延续租约的动作
	for {
		select {
		case <-time.After(5 * time.Second):
			// 续约操作
			etcdClient.KeepAliveOnce(ctx, lease.ID)
		case <-ctx.Done():
			return
		}
	}
}

// func main() {
// 	const EtcdAddr = "http://localhost:2379"
// 	var err error
// 	// 创建 etcd 客户端
// 	etcdClient, err := eclient.NewFromURL(EtcdAddr)
// 	if err != nil {
// 		log.Fatalln(err.Error())
// 		return
// 	}

// 	// 创建 etcd 实现的 grpc 服务注册发现模块 resolver
// 	etcdResolverBuilder, err := eresolver.NewBuilder(etcdClient)
// 	if err != nil {
// 		log.Fatalln(err.Error())
// 		return
// 	}

// 	// 创建 grpc 连接代理
// 	conn, err := grpc.Dial(
// 		// 服务名称
// 		"trade/server1",
// 		// 注入 etcd resolver
// 		grpc.WithResolvers(etcdResolverBuilder),
// 		// 声明使用的负载均衡策略为 roundrobin，轮询。（测试 target 时去除该注释）
// 		// grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"LoadBalancingPolicy": "%s"}`, roundrobin.Name)),
// 		grpc.WithTransportCredentials(insecure.NewCredentials()),
// 	)
// 	if err != nil {
// 		log.Fatalln(err.Error())
// 		return
// 	}

// 	for i := 0; i < 4; i++ {
// 		proxyClient := openTrade.NewOpenServiceClient(conn)

// 		// 通过context携带信息：该请求属于哪个标的
// 		header := metadata.New(map[string]string{"symbol": "s1"})
// 		ctx := metadata.NewOutgoingContext(context.Background(), header)

// 		reply, _ := proxyClient.OpenTrade(ctx, &openTrade.OpenRequest{Symbol: "s1", Price: "100"})
// 		log.Printf("Reply from OpenTrade server: %s", reply.String())
// 	}

// 	defer conn.Close()
// }
