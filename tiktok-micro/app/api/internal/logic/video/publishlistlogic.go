package video

import (
	"context"

	"tiktok-micro/app/api/internal/svc"
	"tiktok-micro/app/api/internal/types"
	"tiktok-micro/app/services/video/video_pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type PublishListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPublishListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PublishListLogic {
	return &PublishListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PublishListLogic) PublishList(req *types.PublishListRequest) (resp *video_pb.PublishlistResponse, err error) {
	res, err := l.svcCtx.VideoRpc.PublishList(
		l.ctx, &video_pb.PublishlistRequest{
			UserId: req.UserID,
			Token:  req.Token,
		})

	if err != nil {
		return nil, err

	}

	return res, nil
}
