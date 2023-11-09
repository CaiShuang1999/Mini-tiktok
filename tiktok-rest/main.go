package main

import (
	"tiktok-rest/route"
	"tiktok-rest/utils"

	"github.com/gin-gonic/gin"
)

func main() {

	utils.InitConfig()
	utils.InitMysql()
	utils.InitRedis()
	utils.GetStaticUrl()

	r := gin.Default() //通过使用gin.Default()，你可以快速设置一个具有常用中间件和全局路由的Gin引擎
	route.RouteInit(r)

	// 设置静态文件路由和目录

	r.Run(":8080")

}
