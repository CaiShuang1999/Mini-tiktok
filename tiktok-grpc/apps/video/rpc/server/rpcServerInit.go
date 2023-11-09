package server

import (
	"fmt"
	"log"
	"net"
	"tiktok-grpc/apps/video/video_pb"

	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func VideoServerInit() {
	serviceHost := viper.GetString("videoService.host")
	servicePort := viper.GetString("videoService.port")
	serviceAddress := serviceHost + ":" + servicePort

	videoRpcServer := grpc.NewServer()
	video_pb.RegisterVideoServiceServer(videoRpcServer, &VideoServiceServer{})

	ls, err := net.Listen("tcp", serviceAddress)
	if err != nil {
		fmt.Println("err:", err.Error())

	}

	if err := videoRpcServer.Serve(ls); err != nil {
		log.Fatalf("视频微服务RPC服务端启动失败: %v", err)

	}

}
