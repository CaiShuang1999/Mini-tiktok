package logic

import (
	"context"

	"tiktok-micro/app/services/user/userservice"
	"tiktok-micro/app/services/video/internal/svc"
	"tiktok-micro/app/services/video/video_pb"
	"tiktok-micro/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type PublishListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPublishListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PublishListLogic {
	return &PublishListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *PublishListLogic) PublishList(in *video_pb.PublishlistRequest) (*video_pb.PublishlistResponse, error) {

	userID := in.UserId
	var videos []model.Video
	db := l.svcCtx.DB

	//db.Preload("Author").Where("user_id = ?", userID).Find(&videos) //使用db.Preload("Author")预加载user表中"Author"信息，否则为0值

	db.Where("user_id = ?", userID).Find(&videos)

	userinfores, _ := l.svcCtx.UserRpc.GetUserInfo(ctx, &userservice.UserInfoRequest{
		UserId: userID,
		Token:  in.Token,
	})

	videos_msg := []*video_pb.VideoInfo{}
	for i := range videos {

		x := &video_pb.VideoInfo{
			Id: videos[i].ID,
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
			CommentCount:  videos[i].CommentCount,
			FavoriteCount: videos[i].FavoriteCount,
			CoverUrl:      videos[i].CoverURL,
			PlayUrl:       videos[i].PlayURL,
			Title:         videos[i].Title,
			CreateTime:    videos[i].CreateTime,
		}
		videos_msg = append(videos_msg, x)
		//videos_msg = append(videos_msg, ConvertVideoToProto(videos[i]))
	}

	return &video_pb.PublishlistResponse{
		StatusCode: 0,
		StatusMsg:  "发布列表",
		VideoList:  videos_msg,
	}, nil

}
func ConvertVideoToProto(video model.Video) *video_pb.VideoInfo {
	protoVideo := &video_pb.VideoInfo{
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

func ConvertAuthorToProto(user model.User) *video_pb.UserInfo {
	protoUser := &video_pb.UserInfo{

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
