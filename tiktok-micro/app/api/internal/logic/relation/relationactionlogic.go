package relation

import (
	"context"

	"tiktok-micro/app/api/internal/svc"
	"tiktok-micro/app/api/internal/types"
	"tiktok-micro/app/services/relation/relation_pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type RelationActionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRelationActionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RelationActionLogic {
	return &RelationActionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RelationActionLogic) RelationAction(req *types.RelationActionRequest) (resp *relation_pb.RelationActionResponse, err error) {
	// todo: add your logic here and delete this line
	res, err := l.svcCtx.RelationRpc.RelationAction(
		l.ctx, &relation_pb.RelationActionRequest{
			ToUserId:   req.ToUserId,
			Token:      req.Token,
			ActionType: req.ActionType,
		})
	if err != nil {
		return nil, err
	}
	return res, err
}
