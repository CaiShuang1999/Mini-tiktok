package logic

import (
	"context"
	"strconv"

	"tiktok-micro/app/services/user/userservice"
	"tiktok-micro/app/services/video/internal/svc"
	"tiktok-micro/app/services/video/video_pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetVideoInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetVideoInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetVideoInfoLogic {
	return &GetVideoInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetVideoInfoLogic) GetVideoInfo(in *video_pb.GetVideoInfoRequest) (*video_pb.GetVideoInfoResponse, error) {
	// todo: add your logic here and delete this line

	video_id := in.VideoId
	video_id_int, _ := strconv.Atoi(video_id)

	video, err := VideoCache(l.svcCtx.DB, l.svcCtx.Rdb_0, l.svcCtx.Rdb_1, int64(video_id_int))
	if err != nil {
		return nil, err
	}
	author_id := strconv.Itoa(int(video.UserID))

	userinfores, _ := l.svcCtx.UserRpc.GetUserInfo(ctx, &userservice.UserInfoRequest{
		UserId: author_id,
		Token:  in.Token,
	})

	return &video_pb.GetVideoInfoResponse{
		StatusCode: 0, StatusMsg: "获取视频信息",
		VideoInfo: &video_pb.VideoInfo{
			Id: int64(video_id_int),
			Author: &video_pb.UserInfo{
				Id:              userinfores.User.Id,
				Name:            userinfores.User.Name,
				Avatar:          userinfores.User.Avatar,
				BackgroundImage: userinfores.User.BackgroundImage,
				FavoriteCount:   userinfores.User.FavoriteCount,
				FollowCount:     userinfores.User.FollowCount,
				FollowerCount:   userinfores.User.FollowerCount,
				WorkCount:       userinfores.User.WorkCount,
				TotalFavorited:  userinfores.User.TotalFavorited,
				Signature:       userinfores.User.Signature,
			},
			CommentCount:  video.CommentCount,
			FavoriteCount: video.FavoriteCount,
			CoverUrl:      video.CoverURL,
			PlayUrl:       video.PlayURL,
			Title:         video.Title,
			CreateTime:    video.CreateTime,
		},
	}, nil
}
