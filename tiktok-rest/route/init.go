package route

import (
	"tiktok-rest/service"

	"github.com/gin-gonic/gin"
)

func RouteInit(r *gin.Engine) {

	apiRouter := r.Group("/douyin")
	apiRouter.GET("/feed/", service.FeedList)

	userGroup := apiRouter.Group("/user")
	{
		userGroup.POST("/register/", service.CreateUser) //注册
		userGroup.POST("/login/", service.LoginUser)     //登录
		userGroup.GET("/", service.UserInfo)             //用户信息
	}

	publishGroup := apiRouter.Group("/publish")
	{
		publishGroup.POST("/action/", service.PublishAPI)  //投稿接口
		publishGroup.GET("/list/", service.PublishListAPI) //用户发布列表
	}

	favoriteGroup := apiRouter.Group("/favorite")
	{
		favoriteGroup.POST("/action/", service.FavoriteAction) //赞操作
		favoriteGroup.GET("/list/", service.FavoriteList)      //喜欢列表
	}

	commentGroup := apiRouter.Group("/comment")
	{
		commentGroup.POST("/action/", service.CommentAction) //评论操作
		commentGroup.GET("/list/", service.CommentList)      //评论列表
	}

	relationGroup := apiRouter.Group("/relation")
	{
		relationGroup.POST("/action/", service.RelationAction) //关注操作
		relationGroup.GET("/follow/list/", service.FollowList) //关注列表
		relationGroup.GET("/follower/list/", service.FansList) //粉丝列表
		relationGroup.GET("/friend/list/", service.FriendList) //好友列表
	}

	messageGroup := apiRouter.Group("/message")
	{
		messageGroup.POST("/action/", service.MessageAction) //发送信息
		messageGroup.GET("/chat/", service.MessageList)      //聊天记录
	}

}
