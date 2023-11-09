package server

import (
	"fmt"
	"log"
	"net"
	"tiktok-grpc/apps/message/messages_pb"

	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func MessageServerInit() {
	serviceHost := viper.GetString("messageService.host")
	servicePort := viper.GetString("messageService.port")
	serviceAddress := serviceHost + ":" + servicePort

	messageRpcServer := grpc.NewServer()
	messages_pb.RegisterMessageServiceServer(messageRpcServer, &MessageServiceServer{})

	ls, err := net.Listen("tcp", serviceAddress)
	if err != nil {
		fmt.Println("err:", err.Error())

	}

	if err := messageRpcServer.Serve(ls); err != nil {
		log.Fatalf("消息微服务RPC服务端启动失败: %v", err)

	}

}
