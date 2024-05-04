package main

import "tiktok-grpc/cmd"

func main() {
	cmd.InitSQL()
	cmd.InitRedis()

	//静态文件服务器
	cmd.GetStaticUrl()

}
