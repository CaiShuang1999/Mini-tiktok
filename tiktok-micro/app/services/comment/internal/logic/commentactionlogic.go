package logic

import (
	"context"
	"strconv"
	"time"

	"tiktok-micro/app/services/comment/comment_pb"
	"tiktok-micro/app/services/comment/internal/svc"
	"tiktok-micro/app/services/user/user_pb"
	"tiktok-micro/common/jwtx"
	"tiktok-micro/model"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type CommentActionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCommentActionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CommentActionLogic {
	return &CommentActionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CommentActionLogic) CommentAction(in *comment_pb.CommentActionRequest) (*comment_pb.CommentActionResponse, error) {
	// todo: add your logic here and delete this line

	token := in.Token
	videoID := in.VideoId
	actionType := in.ActionType

	tokenmsg, ok := jwtx.ParseToken(token)
	if !ok {
		//token无效

		return &comment_pb.CommentActionResponse{
			StatusCode: 1, StatusMsg: "token无效",
		}, nil
	}

	videoIDInt, err := strconv.Atoi(videoID)
	if err != nil {
		return &comment_pb.CommentActionResponse{
			StatusCode: 1, StatusMsg: "无效video_id",
		}, nil
	}
	db := l.svcCtx.DB
	rdb := l.svcCtx.Rdb_0
	var newComment model.Comment

	switch actionType {
	case "1":
		{

			newComment.CommentMsg = in.CommentText
			newComment.UserID = tokenmsg.UserID
			newComment.VideoID = int64(videoIDInt)
			newComment.CreateDate = time.Now().Format("06-01-02 15:04:05")

			video_id_string := strconv.Itoa(int(newComment.VideoID))
			db.Create(&newComment)

			db.Model(&model.Video{}).Where("id = ?", newComment.VideoID).Update("comment_count", gorm.Expr("comment_count  + ?", 1))
			err := rdb.Del(context.Background(), "video:"+video_id_string).Err()
			if err != nil {
				return nil, err
			}
			user_id_string := strconv.Itoa(int(newComment.UserID))
			userinfo, err := l.svcCtx.UserRpc.GetUserInfo(context.Background(), &user_pb.UserInfoRequest{
				UserId: user_id_string,
				Token:  token,
			})
			if err != nil {
				return nil, err
			}
			return &comment_pb.CommentActionResponse{
				StatusCode: 0, StatusMsg: "评论成功",
				Comment: &comment_pb.Comment{
					Id:         newComment.ID,
					Content:    newComment.CommentMsg,
					CreateDate: newComment.CreateDate,
					User: &comment_pb.UserInfo{
						Id:              userinfo.User.Id,
						Name:            userinfo.User.Name,
						FollowCount:     userinfo.User.FollowCount,
						FavoriteCount:   userinfo.User.FavoriteCount,
						FollowerCount:   userinfo.User.FollowerCount,
						Avatar:          userinfo.User.Avatar,
						BackgroundImage: userinfo.User.BackgroundImage,
						TotalFavorited:  userinfo.User.TotalFavorited,
						Signature:       userinfo.User.Signature,
						WorkCount:       userinfo.User.WorkCount,
					},
				},
			}, nil
		}

	case "2":
		{

			commentID := in.CommentId
			commentIDInt, _ := strconv.Atoi(commentID)

			db.First(&newComment, commentID)

			if newComment.UserID != tokenmsg.UserID {
				return &comment_pb.CommentActionResponse{
					StatusCode: 1, StatusMsg: "无法删除别人的评论!",
				}, nil

			}

			if newComment.DeleteDate != "" {
				return &comment_pb.CommentActionResponse{
					StatusCode: 1, StatusMsg: "已经删除!",
				}, nil
			}

			//软删除
			db.Model(&model.Comment{}).Where("id=?", int64(commentIDInt)).
				Update("delete_date", time.Now().Format("06-01-02 15:04:05"))

			db.Model(&model.Video{}).Where("id = ?", int64(videoIDInt)).Update("comment_count", gorm.Expr("comment_count  - ?", 1))

			return &comment_pb.CommentActionResponse{
				StatusCode: 0, StatusMsg: "删除评论成功!",
			}, nil

		}
	default:
		{
			return &comment_pb.CommentActionResponse{
				StatusCode: 1, StatusMsg: "无效请求!",
			}, nil
		}
	}
}
