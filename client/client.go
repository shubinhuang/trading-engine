package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"trading-engine/protos/closeTrade"
	"trading-engine/protos/openTrade"
	"trading-engine/protos/processOrder"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// 开启指定标的的撮合
func openCli(conn *grpc.ClientConn, symbol string, price string) {
	// 创建指定服务的客户端
	openClient := openTrade.NewOpenServiceClient(conn)
	// 通过context携带信息：该请求属于哪个标的
	header := metadata.New(map[string]string{"symbol": symbol})
	ctx := metadata.NewOutgoingContext(context.Background(), header)

	// grpc远程调用服务中定义的方法
	reply, err := openClient.OpenTrade(ctx, &openTrade.OpenRequest{Symbol: symbol, Price: price})
	if err != nil {
		log.Printf("Error when calling OpenTrade: %s", err)
	}
	log.Printf("Reply from OpenTrade server: %s", reply.String())
}

// 关闭指定标的的撮合
func closeCli(conn *grpc.ClientConn, symbol string) {
	// 创建指定服务的客户端
	closeClient := closeTrade.NewCloseServiceClient(conn)
	// 通过context携带信息：该请求属于哪个标的
	header := metadata.New(map[string]string{"symbol": symbol})
	ctx := metadata.NewOutgoingContext(context.Background(), header)

	// grpc远程调用服务中定义的方法
	reply, err := closeClient.CloseTrade(ctx, &closeTrade.CloseRequest{Symbol: symbol})
	if err != nil {
		log.Printf("Error when calling CloseTrade: %s", err)
	}
	log.Printf("Reply from CloseTrade server: %s", reply.String())
}

// 下单测试
func createOrder(conn *grpc.ClientConn, symbol string, userId string, orderId string, direction string, price string, quantity string) {
	// 创建指定服务的客户端
	orderClient := processOrder.NewOrderServiceClient(conn)
	// 通过context携带信息：该请求属于哪个标的
	header := metadata.New(map[string]string{"symbol": symbol})
	ctx := metadata.NewOutgoingContext(context.Background(), header)

	// grpc远程调用服务中定义的方法
	reply, err := orderClient.CreateOrder(
		ctx,
		&processOrder.CreateRequest{
			Symbol:    symbol,
			UserId:    userId,
			OrderId:   orderId,
			Direction: direction,
			Price:     price,
			Quantity:  quantity})
	if err != nil {
		log.Printf("Error when calling CreateOrder: %s", err)
	}
	log.Printf("Reply from CreateOrder server: %s", reply.String())
}

// 撤单测试
func cancelOrder(conn *grpc.ClientConn, symbol string, orderId string, direction string) {
	// 创建指定服务的客户端
	orderClient := processOrder.NewOrderServiceClient(conn)
	// 通过context携带信息：该请求属于哪个标的
	header := metadata.New(map[string]string{"symbol": symbol})
	ctx := metadata.NewOutgoingContext(context.Background(), header)

	// grpc远程调用服务中定义的方法
	reply, err := orderClient.CancelOrder(ctx, &processOrder.CancelRequest{Symbol: symbol, OrderId: orderId, Direction: direction})
	if err != nil {
		log.Printf("Error when calling CancelOrder: %s", err)
	}
	log.Printf("Reply from CancelOrder server: %s", reply.String())
}

func main() {
	// consul.InitConsul()
	// services := consul.DiscoveryService("proxy", "")
	// var addr = services[0].Node.Address + ":" + strconv.Itoa(services[0].Service.Port)
	// // 连接grpc服务器
	// conn, err := grpc.Dial(addr, grpc.WithInsecure(), grpc.WithBlock())
	conn, err := grpc.Dial("localhost:5432", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("Did not connect: %v", err)
	}
	defer conn.Close()

	fmt.Printf("Trade engine client\nInput help for usages\n\n$ ")
	reader := bufio.NewReader(os.Stdin)

	for {
		text, _ := reader.ReadString('\n')
		text = strings.TrimSpace(text)
		args := strings.Split(text, " ")

		if args[0] == "opentrade" {
			// opentrade symbol price
			if len(args) == 3 {
				openCli(conn, args[1], args[2])
			}
		} else if args[0] == "closetrade" {
			// closetrade symbol
			if len(args) == 2 {
				closeCli(conn, args[1])
			}
		} else if args[0] == "createorder" {
			// createorder symbol userId orderId direction price quantity
			if len(args) == 7 {
				createOrder(conn, args[1], args[2], args[3], args[4], args[5], args[6])
			}
		} else if args[0] == "cancelorder" {
			// cancelorder symbol orderId direction
			if len(args) == 4 {
				cancelOrder(conn, args[1], args[2], args[3])
			}
		} else if args[0] == "help" {
			fmt.Printf("Usages:\n----------\n")
			fmt.Println("开启指定标的的撮合: opentrade symbol price")
			fmt.Println("关闭指定标的的撮合: closetrade symbol")
			fmt.Println("下单: createorder symbol userId orderId direction price quantity")
			fmt.Println("撤单: cancelorder symbol orderId direction")
			fmt.Println("帮助: help")
			fmt.Println("退出: exit")
		} else if args[0] == "exit" {
			break
		} else {
			fmt.Println("wrong command")
		}

		fmt.Printf("\n$ ")
	}
}
