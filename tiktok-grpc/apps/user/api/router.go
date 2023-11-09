package api

import (
	"tiktok-grpc/apps/user/rpc/client"

	"github.com/gin-gonic/gin"
)

func InitUserRouter(r *gin.Engine) {
	apiRouter := r.Group("/douyin")

	userGroup := apiRouter.Group("/user")
	{
		userGroup.POST("/register/", client.Register) //注册
		userGroup.POST("/login/", client.Login)       //登录
		userGroup.GET("/", client.GetUserInfo)        //用户信息
	}

}
