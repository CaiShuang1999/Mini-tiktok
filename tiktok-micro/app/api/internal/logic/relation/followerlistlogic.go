package relation

import (
	"context"

	"tiktok-micro/app/api/internal/svc"
	"tiktok-micro/app/api/internal/types"
	"tiktok-micro/app/services/relation/relation_pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type FollowerListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFollowerListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FollowerListLogic {
	return &FollowerListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FollowerListLogic) FollowerList(req *types.FollowerListRequest) (resp *relation_pb.GetFansListResponse, err error) {
	res, err := l.svcCtx.RelationRpc.GetFansList(l.ctx, &relation_pb.GetFansListRequest{
		UserId: req.UserId,
		Token:  req.Token,
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}
