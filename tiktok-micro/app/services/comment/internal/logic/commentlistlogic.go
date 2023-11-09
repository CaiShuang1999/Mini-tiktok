package logic

import (
	"context"

	"tiktok-micro/app/services/comment/comment_pb"
	"tiktok-micro/app/services/comment/internal/svc"
	"tiktok-micro/common/jwtx"
	"tiktok-micro/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type CommentListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCommentListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CommentListLogic {
	return &CommentListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CommentListLogic) CommentList(in *comment_pb.CommentListRequest) (*comment_pb.CommentListResponse, error) {
	token := in.Token
	videoID := in.VideoId

	_, ok := jwtx.ParseToken(token)
	if !ok {
		//token无效
		return &comment_pb.CommentListResponse{
			StatusCode: 1, StatusMsg: "无效token!",
		}, nil
	}

	db := l.svcCtx.DB
	var comments []model.Comment

	db.Preload("UserMsg").Where("delete_date =? AND video_id=?", "", videoID).Order("id desc").Find(&comments)

	comments_msg := []*comment_pb.Comment{}
	for i := range comments {
		comments_msg = append(comments_msg,
			&comment_pb.Comment{
				Id:         comments[i].ID,
				Content:    comments[i].CommentMsg,
				CreateDate: comments[i].CreateDate,
				User: &comment_pb.UserInfo{
					Id:              comments[i].UserID,
					Name:            comments[i].UserMsg.Name,
					Avatar:          comments[i].UserMsg.Avatar,
					BackgroundImage: comments[i].UserMsg.BackgroundImage,
					FavoriteCount:   comments[i].UserMsg.FavoriteCount,
					FollowCount:     comments[i].UserMsg.FollowCount,
					FollowerCount:   comments[i].UserMsg.FollowerCount,
					WorkCount:       comments[i].UserMsg.WorkCount,
					Signature:       comments[i].UserMsg.Signature,
					TotalFavorited:  comments[i].UserMsg.TotalFavorited,
				},
			})
	}

	return &comment_pb.CommentListResponse{
		StatusCode: 0, StatusMsg: "评论区",
		CommentList: comments_msg,
	}, nil
}
