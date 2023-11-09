package server

import (
	"fmt"
	"log"
	"net"
	"tiktok-grpc/apps/comment/comment_pb"

	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func CommentServerInit() {
	serviceHost := viper.GetString("commentService.host")
	servicePort := viper.GetString("commentService.port")
	serviceAddress := serviceHost + ":" + servicePort

	commentRpcServer := grpc.NewServer()
	comment_pb.RegisterCommentServiceServer(commentRpcServer, &commentServiceServer{})
	ls, err := net.Listen("tcp", serviceAddress)
	if err != nil {
		fmt.Println("err:", err.Error())

	}

	if err := commentRpcServer.Serve(ls); err != nil {
		log.Fatalf("评论微服务RPC服务端启动失败: %v", err)

	}

}
