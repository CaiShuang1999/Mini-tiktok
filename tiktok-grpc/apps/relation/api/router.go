package api

import (
	"tiktok-grpc/apps/relation/rpc/client"

	"github.com/gin-gonic/gin"
)

func InitRelationRouter(r *gin.Engine) {
	apiRouter := r.Group("/douyin")

	relationGroup := apiRouter.Group("/relation")
	{
		relationGroup.POST("/action/", client.RelationAction)    //关注操作
		relationGroup.GET("/follow/list/", client.GetFollowList) //关注列表
		relationGroup.GET("/follower/list/", client.GetFanList)  //粉丝列表
		relationGroup.GET("/friend/list/", client.GetFriendList) //好友列表
	}

}
