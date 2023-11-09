package api

import (
	"tiktok-grpc/apps/feed/rpc/client"

	"github.com/gin-gonic/gin"
)

func InitUserRouter(r *gin.Engine) {
	apiRouter := r.Group("/douyin")

	apiRouter.GET("/feed/", client.FeedList)

}
