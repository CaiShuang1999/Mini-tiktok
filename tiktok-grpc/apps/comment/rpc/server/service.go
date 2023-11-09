package server

import (
	"context"
	"strconv"
	"tiktok-grpc/apps/comment/comment_pb"
	"tiktok-grpc/cmd"
	"tiktok-grpc/common/jwtx"
	"tiktok-grpc/common/utils"
	"tiktok-grpc/model"
	"time"

	"gorm.io/gorm"
)

type commentServiceServer struct {
	comment_pb.UnimplementedCommentServiceServer
}

func (p *commentServiceServer) CommentAction(ctx context.Context, req *comment_pb.CommentActionRequest) (*comment_pb.CommentActionResponse, error) {
	token := req.Token
	videoID := req.VideoId
	actionType := req.ActionType

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
	db := cmd.DB

	var newComment model.Comment

	switch actionType {
	case "1":
		{

			newComment.CommentMsg = req.CommentText
			newComment.UserID = tokenmsg.UserID
			newComment.VideoID = int64(videoIDInt)
			newComment.CreateDate = time.Now().Format("06-01-02 15:04:05")

			db.Create(&newComment)
			db.First(&newComment.UserMsg, "id = ?", newComment.UserID)
			db.Model(&model.Video{}).Where("id = ?", newComment.VideoID).Update("comment_count", gorm.Expr("comment_count  + ?", 1))

			return &comment_pb.CommentActionResponse{
				StatusCode: 0, StatusMsg: "评论成功",
				Comment: utils.ConvertCommentToProto(newComment),
			}, nil
		}

	case "2":
		{

			commentID := req.CommentId
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
func (p *commentServiceServer) CommentList(ctx context.Context, req *comment_pb.CommentListRequest) (*comment_pb.CommentListResponse, error) {
	token := req.Token
	videoID := req.VideoId

	_, ok := jwtx.ParseToken(token)
	if !ok {
		//token无效
		return &comment_pb.CommentListResponse{
			StatusCode: 1, StatusMsg: "无效token!",
		}, nil
	}

	db := cmd.DB
	var comments []model.Comment

	db.Preload("UserMsg").Where("delete_date =? AND video_id=?", "", videoID).Order("id desc").Find(&comments)

	comments_msg := []*comment_pb.Comment{}
	for i := range comments {
		comments_msg = append(comments_msg, utils.ConvertCommentToProto(comments[i]))
	}

	return &comment_pb.CommentListResponse{
		StatusCode: 0, StatusMsg: "评论区",
		CommentList: comments_msg,
	}, nil
}
