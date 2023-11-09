package cmd

import (
	comment_api "tiktok-grpc/apps/comment/api"
	favorite_api "tiktok-grpc/apps/favorite/api"
	feed_api "tiktok-grpc/apps/feed/api"
	message_api "tiktok-grpc/apps/message/api"
	relation_api "tiktok-grpc/apps/relation/api"
	user_api "tiktok-grpc/apps/user/api"
	video_api "tiktok-grpc/apps/video/api"

	"github.com/gin-gonic/gin"
)

func InitRouter(r *gin.Engine) {

	user_api.InitUserRouter(r)
	video_api.InitVideoRouter(r)
	feed_api.InitUserRouter(r)
	favorite_api.InitFavoriteRouter(r)
	comment_api.InitCommentRouter(r)
	relation_api.InitRelationRouter(r)
	message_api.InitMessageRouter(r)

}
