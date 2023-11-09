package api

import (
	"tiktok-grpc/apps/video/rpc/client"

	"github.com/gin-gonic/gin"
)

func InitVideoRouter(r *gin.Engine) {
	apiRouter := r.Group("/douyin")

	publishGroup := apiRouter.Group("/publish")
	{
		publishGroup.POST("/action/", client.UploadVideo) //投稿接口
		publishGroup.GET("/list/", client.PubilshList)    //用户发布列表
	}

}
