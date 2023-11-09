package api

import (
	"tiktok-grpc/apps/favorite/rpc/client"

	"github.com/gin-gonic/gin"
)

func InitFavoriteRouter(r *gin.Engine) {
	apiRouter := r.Group("/douyin")

	favoriteGroup := apiRouter.Group("/favorite")
	{
		favoriteGroup.POST("/action/", client.FavoriteAction) //赞操作
		favoriteGroup.GET("/list/", client.FavoriteList)      //喜欢列表
	}

}
