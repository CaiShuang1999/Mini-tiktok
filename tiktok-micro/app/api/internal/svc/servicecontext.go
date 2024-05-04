package svc

import (
	"tiktok-micro/app/api/internal/config"
	"tiktok-micro/app/services/comment/commentservice"
	"tiktok-micro/app/services/favorite/favoriteservice"
	"tiktok-micro/app/services/feed/feedservice"
	"tiktok-micro/app/services/message/messageservice"
	"tiktok-micro/app/services/relation/relationservice"
	"tiktok-micro/app/services/user/userservice"
	"tiktok-micro/app/services/video/videoservice"

	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config   config.Config
	UserRpc  userservice.UserService
	VideoRpc videoservice.VideoService
	FeedRpc  feedservice.FeedService

	FavoriteRpc favoriteservice.FavoriteService
	CommentRpc  commentservice.CommentService
	RelationRpc relationservice.RelationService
	MessageRpc  messageservice.MessageService
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:   c,
		UserRpc:  userservice.NewUserService(zrpc.MustNewClient(c.UserRpc)),
		VideoRpc: videoservice.NewVideoService(zrpc.MustNewClient(c.VideoRpc)),
		FeedRpc:  feedservice.NewFeedService(zrpc.MustNewClient(c.FeedRpc)),

		FavoriteRpc: favoriteservice.NewFavoriteService(zrpc.MustNewClient(c.FavoriteRpc)),
		CommentRpc:  commentservice.NewCommentService(zrpc.MustNewClient(c.CommentRpc)),
		RelationRpc: relationservice.NewRelationService(zrpc.MustNewClient(c.RelationRpc)),
		MessageRpc:  messageservice.NewMessageService(zrpc.MustNewClient(c.MessageRpc)),
	}
}
