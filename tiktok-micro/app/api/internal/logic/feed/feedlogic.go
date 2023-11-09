package feed

import (
	"context"

	"tiktok-micro/app/api/internal/svc"
	"tiktok-micro/app/api/internal/types"
	"tiktok-micro/app/services/feed/feed_pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type FeedLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFeedLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FeedLogic {
	return &FeedLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FeedLogic) Feed(req *types.FeedRequest) (resp *feed_pb.FeedResponse, err error) {
	res, err := l.svcCtx.FeedRpc.FeedList(
		l.ctx, &feed_pb.FeedRequest{
			Token:      req.Token,
			LatestTime: req.LatestTime,
		})
	if err != nil {
		return nil, err
	}

	return res, nil

}
