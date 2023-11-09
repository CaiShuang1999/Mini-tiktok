package api

import (
	"tiktok-grpc/apps/comment/rpc/client"

	"github.com/gin-gonic/gin"
)

func InitCommentRouter(r *gin.Engine) {
	apiRouter := r.Group("/douyin")

	commentGroup := apiRouter.Group("/comment")
	{
		commentGroup.POST("/action/", client.CommentAction) //评论操作
		commentGroup.GET("/list/", client.CommentList)      //评论列表
	}

}
