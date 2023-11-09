package comment

import (
	"context"

	"tiktok-micro/app/api/internal/svc"
	"tiktok-micro/app/api/internal/types"
	"tiktok-micro/app/services/comment/comment_pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type CommentListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCommentListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CommentListLogic {
	return &CommentListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CommentListLogic) CommentList(req *types.CommentListRequest) (resp *comment_pb.CommentListResponse, err error) {
	res, err := l.svcCtx.CommentRpc.CommentList(
		l.ctx, &comment_pb.CommentListRequest{
			Token:   req.Token,
			VideoId: req.VideoId,
		})
	if err != nil {
		return nil, err
	}

	return res, nil
}
