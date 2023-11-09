package logic

import (
	"context"

	"tiktok-micro/app/services/feed/feed_pb"
	"tiktok-micro/app/services/feed/internal/svc"
	"tiktok-micro/common/jwtx"
	"tiktok-micro/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type FeedListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFeedListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FeedListLogic {
	return &FeedListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FeedListLogic) FeedList(in *feed_pb.FeedRequest) (*feed_pb.FeedResponse, error) {
	// todo: add your logic here and delete this line
	token := in.Token

	var videos []model.Video

	db := l.svcCtx.DB
	db.Preload("Author").Order("id desc").Limit(30).Find(&videos)

	if token != "" {
		tokenmsg, ok := jwtx.ParseToken(token)
		if !ok {
			return &feed_pb.FeedResponse{
				StatusCode: 1,
				StatusMsg:  "token无效",
			}, nil
		}
		userID := tokenmsg.UserID

		for i := range videos {
			if db.Where("video_id = ? AND user_id = ? AND is_favorite=?", videos[i].ID, userID, true).Find(&model.Favorite{}).RowsAffected != 0 {
				videos[i].IsFavorite = true
			}

			if db.Where("user_id =? AND to_user_id=? AND is_follow=?", userID, videos[i].Author.ID, true).Find(&model.Relation{}).RowsAffected != 0 {
				videos[i].Author.IsFollow = true
			}

		}

	}
	for i := range videos {
		videos[i].Author.Avatar = "http://" + l.svcCtx.Config.Nginx.Addr + videos[i].Author.Avatar
		videos[i].Author.BackgroundImage = "http://" + l.svcCtx.Config.Nginx.Addr + videos[i].Author.BackgroundImage
	}
	videos_msg := []*feed_pb.VideoInfo{}
	for i := range videos {
		videos_msg = append(videos_msg, ConvertVideoToProto(videos[i]))
	}

	return &feed_pb.FeedResponse{
		StatusCode: 0,
		StatusMsg:  "feed",
		VideoList:  videos_msg,
	}, nil
}
func ConvertVideoToProto(video model.Video) *feed_pb.VideoInfo {
	protoVideo := &feed_pb.VideoInfo{
		Id:            video.ID,
		Author:        ConvertAuthorToProto(video.Author),
		CommentCount:  video.CommentCount,
		CoverUrl:      video.CoverURL,
		PlayUrl:       video.PlayURL,
		Title:         video.Title,
		CreateTime:    video.CreateTime,
		FavoriteCount: video.FavoriteCount,
		IsFavorite:    video.IsFavorite,
	}

	return protoVideo
}

func ConvertAuthorToProto(user model.User) *feed_pb.UserInfo {
	protoUser := &feed_pb.UserInfo{

		Id:              user.ID,
		Name:            user.Name,
		FollowCount:     user.FollowCount,
		FollowerCount:   user.FollowerCount,
		Avatar:          user.Avatar,
		BackgroundImage: user.BackgroundImage,
		Signature:       user.Signature,
		TotalFavorited:  user.TotalFavorited,
		WorkCount:       user.WorkCount,
		FavoriteCount:   user.FavoriteCount,
	}

	return protoUser
}
