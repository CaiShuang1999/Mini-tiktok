package relation

import (
	"context"

	"tiktok-micro/app/api/internal/svc"
	"tiktok-micro/app/api/internal/types"
	"tiktok-micro/app/services/relation/relation_pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type FollowListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFollowListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FollowListLogic {
	return &FollowListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FollowListLogic) FollowList(req *types.FollowListRequest) (resp *relation_pb.GetFollowListResponse, err error) {
	res, err := l.svcCtx.RelationRpc.GetFollowList(
		l.ctx, &relation_pb.GetFollowListRequest{
			UserId: req.UserId,
			Token:  req.Token,
		})

	if err != nil {
		return nil, err
	}

	return res, nil
}
