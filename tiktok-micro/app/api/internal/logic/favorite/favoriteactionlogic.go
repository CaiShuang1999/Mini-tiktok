package favorite

import (
	"context"

	"tiktok-micro/app/api/internal/svc"
	"tiktok-micro/app/api/internal/types"
	"tiktok-micro/app/services/favorite/favorite_pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type FavoriteActionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFavoriteActionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FavoriteActionLogic {
	return &FavoriteActionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FavoriteActionLogic) FavoriteAction(req *types.FavoriteActionRequest) (resp *favorite_pb.FavoriteActionResponse, err error) {
	res, err := l.svcCtx.FavoriteRpc.FavoriteAction(
		l.ctx, &favorite_pb.FavoriteActionRequest{
			Token:      req.Token,
			VideoId:    req.VideoId,
			ActionType: req.ActionType,
		})
	if err != nil {
		return nil, err
	}

	return res, nil
}
