package user

import (
	"context"

	"tiktok-micro/app/api/internal/svc"
	"tiktok-micro/app/api/internal/types"
	"tiktok-micro/app/services/user/user_pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserInfoLogic {
	return &UserInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserInfoLogic) UserInfo(req *types.UserInfoRequest) (resp *types.UserInfoResponse, err error) {
	res, err := l.svcCtx.UserRpc.GetUserInfo(
		l.ctx, &user_pb.UserInfoRequest{
			UserId: req.UserId,
			Token:  req.Token,
		})
	if err != nil {
		return nil, err
	}

	return &types.UserInfoResponse{
		StatusCode: res.StatusCode,
		StatusMsg:  res.StatusMsg,
		User: types.UserInfo{
			ID:              res.User.Id,
			Name:            res.User.Name,
			Avatar:          res.User.Avatar,
			BackgroundImage: res.User.BackgroundImage,
			FavoriteCount:   res.User.FavoriteCount,
			WorkCount:       res.User.WorkCount,
			FollowCount:     res.User.FollowCount,
			FollowerCount:   res.User.FollowerCount,
			TotalFavorited:  res.User.TotalFavorited,
			Signature:       res.User.Signature,
			IsFollow:        res.User.IsFollow,
		},
	}, nil
}
