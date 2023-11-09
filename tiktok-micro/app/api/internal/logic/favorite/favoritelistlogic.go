package favorite

import (
	"context"

	"tiktok-micro/app/api/internal/svc"
	"tiktok-micro/app/api/internal/types"
	"tiktok-micro/app/services/favorite/favorite_pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type FavoriteListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFavoriteListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FavoriteListLogic {
	return &FavoriteListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FavoriteListLogic) FavoriteList(req *types.FavoriteListRequest) (resp *favorite_pb.FavoriteListResponse, err error) {
	res, err := l.svcCtx.FavoriteRpc.FavoriteList(
		l.ctx, &favorite_pb.FavoriteListRequest{
			Token:  req.Token,
			UserId: req.UserId,
		})
	if err != nil {
		return nil, err
	}

	return res, nil
}
