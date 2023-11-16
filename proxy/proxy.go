package main

import (
	"context"
	"encoding/json"
	"hash/fnv"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strconv"
	"trading-engine/consul"

	"github.com/mwitkow/grpc-proxy/proxy"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func hash(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}

// rpc请求发给代理服务器，代理服务器通过context获得请求属于哪个标的，将该请求转发到对应的服务器上
func main() {
	consul.InitConsul()
	// services.json 开启了的撮合服务 ["trade0","trade1"]
	serciceSlice := make([]string, 0)

	jsonFile, err := os.Open("../config/services.json")
	// 如果文件打开失败，需要进行err的错误处理
	if err != nil {
		log.Fatalf("open service.json fail: %v", err)
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &serciceSlice)

	var SERVERNUM uint32 = uint32(len(serciceSlice))

	connPool := make([]*grpc.ClientConn, SERVERNUM)
	for i := 0; i < int(SERVERNUM); i++ {
		services := consul.DiscoveryService(serciceSlice[i], "")
		if len(services) > 0 {
			var addr = services[0].Node.Address + ":" + strconv.Itoa(services[0].Service.Port)
			// 转发给标的对应的地址
			connPool[i], _ = grpc.Dial(addr, grpc.WithInsecure(), grpc.WithCodec(proxy.Codec()))
		}
	}

	// 请求转发规则：对symbol求hash后取模，转发到相应的server
	directorFn := func(ctx context.Context, fullMethodName string) (context.Context, *grpc.ClientConn, error) {
		// log.Println(fullMethodName)
		md, _ := metadata.FromIncomingContext(ctx)
		// log.Printf("%v\n", md)              // map[:authority:[localhost:6543] content-type:[application/grpc] symbol:[s1] user-agent:[grpc-go/1.56.2]]
		// log.Printf("%v\n", md["symbol"][0]) // s1

		if val, exists := md["symbol"]; exists {
			conn := connPool[hash(val[0])%SERVERNUM]
			return ctx, conn, nil
		}
		return ctx, nil, status.Errorf(codes.Unimplemented, "Unknown method")
	}

	// 设置并运行代理服务器
	srv := grpc.NewServer(
		grpc.CustomCodec(proxy.Codec()),
		grpc.UnknownServiceHandler(proxy.TransparentHandler(directorFn)))

	// 暴露6543端口，提供proxy服务
	lis, err := net.Listen("tcp", ":6543")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Printf("proxy is listening...")

	if err := srv.Serve(lis); err != nil {
		log.Fatalf("Fail to serve: %v", err)
	}
}
