package server

import (
	"fmt"
	"log"
	"net"
	"tiktok-grpc/apps/feed/feed_pb"

	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func FeedServerInit() {
	serviceHost := viper.GetString("feedService.host")
	servicePort := viper.GetString("feedService.port")
	serviceAddress := serviceHost + ":" + servicePort

	feedRpcServer := grpc.NewServer()
	feed_pb.RegisterFeedServiceServer(feedRpcServer, &FeedServiceServer{})

	ls, err := net.Listen("tcp", serviceAddress)
	if err != nil {
		fmt.Println("err:", err.Error())

	}

	if err := feedRpcServer.Serve(ls); err != nil {
		log.Fatalf("feed微服务RPC服务端启动失败: %v", err)

	}

}
