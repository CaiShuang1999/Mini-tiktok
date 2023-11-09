package logic

import (
	"context"
	"strconv"

	"tiktok-micro/app/services/favorite/favorite_pb"
	"tiktok-micro/app/services/favorite/internal/svc"
	"tiktok-micro/app/services/video/video_pb"
	"tiktok-micro/common/jwtx"
	"tiktok-micro/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type FavoriteListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFavoriteListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FavoriteListLogic {
	return &FavoriteListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FavoriteListLogic) FavoriteList(in *favorite_pb.FavoriteListRequest) (*favorite_pb.FavoriteListResponse, error) {
	userID := in.UserId
	userIDint, _ := strconv.Atoi(userID)
	token := in.Token

	tokenmsg, ok := jwtx.ParseToken(token)
	if tokenmsg.UserID != int64(userIDint) && !ok {
		//token无效

		return &favorite_pb.FavoriteListResponse{
			StatusCode: 1, StatusMsg: "无效token",
		}, nil
	}

	db := l.svcCtx.DB
	var favorites []model.Favorite

	db.Where("user_id = ? AND is_favorite=?", userID, true).Find(&favorites)

	videos_msg := []*favorite_pb.VideoInfo{}
	for i := range favorites {
		video_id := strconv.Itoa(int(favorites[i].VideoId))
		res, _ := l.svcCtx.VideoRpc.GetVideoInfo(context.Background(), &video_pb.GetVideoInfoRequest{
			VideoId: video_id,
			Token:   in.Token,
		})
		videos_msg = append(videos_msg, &favorite_pb.VideoInfo{
			Id:            res.VideoInfo.Id,
			CommentCount:  res.VideoInfo.CommentCount,
			FavoriteCount: res.VideoInfo.FavoriteCount,
			CoverUrl:      res.VideoInfo.CoverUrl,
			PlayUrl:       res.VideoInfo.PlayUrl,
			Title:         res.VideoInfo.Title,
			CreateTime:    res.VideoInfo.CreateTime,
			IsFavorite:    true,
			Author: &favorite_pb.UserInfo{
				Id:              res.VideoInfo.Author.Id,
				Name:            res.VideoInfo.Author.Name,
				Avatar:          res.VideoInfo.Author.Avatar,
				BackgroundImage: res.VideoInfo.Author.BackgroundImage,
				FollowCount:     res.VideoInfo.Author.FollowCount,
				FollowerCount:   res.VideoInfo.Author.FollowerCount,
				WorkCount:       res.VideoInfo.Author.WorkCount,
				Signature:       res.VideoInfo.Author.Signature,
				TotalFavorited:  res.VideoInfo.Author.TotalFavorited,
				FavoriteCount:   res.VideoInfo.Author.FavoriteCount,
			},
		})
	}

	return &favorite_pb.FavoriteListResponse{
		StatusCode: 0, StatusMsg: "点赞列表",
		VideoList: videos_msg,
	}, nil
}
