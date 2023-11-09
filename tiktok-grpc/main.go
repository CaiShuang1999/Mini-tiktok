package main

import (
	"tiktok-grpc/apps"
	"tiktok-grpc/cmd"

	"github.com/gin-gonic/gin"
)

func main() {
	//数据库连接
	cmd.InitSQL()
	cmd.InitRedis()

	//静态文件服务器
	cmd.GetStaticUrl()
	//微服务RPC服务端启动
	apps.ServerInits()

	r := gin.Default()
	cmd.InitRouter(r)
	r.Run(":8080")
}
