package comment

import (
	"context"

	"tiktok-micro/app/api/internal/svc"
	"tiktok-micro/app/api/internal/types"
	"tiktok-micro/app/services/comment/comment_pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type CommentActionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCommentActionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CommentActionLogic {
	return &CommentActionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CommentActionLogic) CommentAction(req *types.CommentActionRequest) (resp *comment_pb.CommentActionResponse, err error) {
	res, err := l.svcCtx.CommentRpc.CommentAction(
		l.ctx, &comment_pb.CommentActionRequest{
			Token:       req.Token,
			VideoId:     req.VideoId,
			ActionType:  req.ActionType,
			CommentText: req.CommentText,
			CommentId:   req.CommentID,
		})
	if err != nil {
		return nil, err
	}

	return res, nil
}
