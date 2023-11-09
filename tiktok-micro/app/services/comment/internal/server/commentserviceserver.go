// Code generated by goctl. DO NOT EDIT.
// Source: comment.proto

package server

import (
	"context"

	"tiktok-micro/app/services/comment/comment_pb"
	"tiktok-micro/app/services/comment/internal/logic"
	"tiktok-micro/app/services/comment/internal/svc"
)

type CommentServiceServer struct {
	svcCtx *svc.ServiceContext
	comment_pb.UnimplementedCommentServiceServer
}

func NewCommentServiceServer(svcCtx *svc.ServiceContext) *CommentServiceServer {
	return &CommentServiceServer{
		svcCtx: svcCtx,
	}
}

func (s *CommentServiceServer) CommentAction(ctx context.Context, in *comment_pb.CommentActionRequest) (*comment_pb.CommentActionResponse, error) {
	l := logic.NewCommentActionLogic(ctx, s.svcCtx)
	return l.CommentAction(in)
}

func (s *CommentServiceServer) CommentList(ctx context.Context, in *comment_pb.CommentListRequest) (*comment_pb.CommentListResponse, error) {
	l := logic.NewCommentListLogic(ctx, s.svcCtx)
	return l.CommentList(in)
}