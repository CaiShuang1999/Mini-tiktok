package server

import (
	"fmt"
	"log"
	"net"
	"tiktok-grpc/apps/user/pb"

	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func UserServerInit() {
	serviceHost := viper.GetString("userService.host")
	servicePort := viper.GetString("userService.port")
	serviceAddress := serviceHost + ":" + servicePort

	userRpcServer := grpc.NewServer()
	pb.RegisterUserServiceServer(userRpcServer, &UserServiceServer{})

	ls, err := net.Listen("tcp", serviceAddress)
	if err != nil {
		fmt.Println("err:", err.Error())

	}

	if err := userRpcServer.Serve(ls); err != nil {
		log.Fatalf("用户微服务RPC服务端启动失败: %v", err)

	}

}
