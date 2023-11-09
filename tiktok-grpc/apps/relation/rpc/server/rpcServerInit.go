package server

import (
	"fmt"
	"log"
	"net"
	"tiktok-grpc/apps/relation/relation_pb"

	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func RelationServerInit() {
	serviceHost := viper.GetString("relationService.host")
	servicePort := viper.GetString("relationService.port")
	serviceAddress := serviceHost + ":" + servicePort

	relationRpcServer := grpc.NewServer()
	relation_pb.RegisterRelationServiceServer(relationRpcServer, &RelationServiceServer{})

	ls, err := net.Listen("tcp", serviceAddress)
	if err != nil {
		fmt.Println("err:", err.Error())

	}

	if err := relationRpcServer.Serve(ls); err != nil {
		log.Fatalf("关系微服务RPC服务端启动失败: %v", err)

	}

}
