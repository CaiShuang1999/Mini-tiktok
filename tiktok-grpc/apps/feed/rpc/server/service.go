package server

import (
	"context"
	"tiktok-grpc/apps/feed/feed_pb"
	"tiktok-grpc/cmd"
	"tiktok-grpc/common/jwtx"
	"tiktok-grpc/common/utils"
	"tiktok-grpc/model"
)

type FeedServiceServer struct {
	feed_pb.UnimplementedFeedServiceServer
}

func (p *FeedServiceServer) FeedList(ctx context.Context, req *feed_pb.FeedRequest) (*feed_pb.FeedResponse, error) {

	token := req.Token

	var videos []model.Video

	db := cmd.DB
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
		videos[i].Author.Avatar = cmd.StaticUrl + videos[i].Author.Avatar
		videos[i].Author.BackgroundImage = cmd.StaticUrl + videos[i].Author.BackgroundImage
	}
	videos_msg := []*feed_pb.VideoInfo{}
	for i := range videos {
		videos_msg = append(videos_msg, utils.ConvertVideoTofeedProto(videos[i]))
	}

	return &feed_pb.FeedResponse{
		StatusCode: 0,
		StatusMsg:  "feed",
		VideoList:  videos_msg,
	}, nil
}
