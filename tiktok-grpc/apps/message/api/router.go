package api

import (
	"tiktok-grpc/apps/message/rpc/client"

	"github.com/gin-gonic/gin"
)

func InitMessageRouter(r *gin.Engine) {
	apiRouter := r.Group("/douyin")

	messageGroup := apiRouter.Group("/message")
	{
		messageGroup.POST("/action/", client.MessageAction) //发送信息
		messageGroup.GET("/chat/", client.GetMessageList)   //聊天记录
	}

}
