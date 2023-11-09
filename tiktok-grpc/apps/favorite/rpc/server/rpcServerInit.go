package server

import (
	"fmt"
	"log"
	"net"
	"tiktok-grpc/apps/favorite/favorite_pb"

	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func FavoriteServerInit() {
	serviceHost := viper.GetString("favoriteService.host")
	servicePort := viper.GetString("favoriteService.port")
	serviceAddress := serviceHost + ":" + servicePort

	favoriteRpcServer := grpc.NewServer()
	favorite_pb.RegisterFavoriteServiceServer(favoriteRpcServer, &favoriteServiceServer{})

	ls, err := net.Listen("tcp", serviceAddress)
	if err != nil {
		fmt.Println("err:", err.Error())

	}

	if err := favoriteRpcServer.Serve(ls); err != nil {
		log.Fatalf("点赞微服务RPC服务端启动失败: %v", err)

	}

}
