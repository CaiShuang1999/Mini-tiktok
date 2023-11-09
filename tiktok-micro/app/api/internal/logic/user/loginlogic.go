package user

import (
	"context"

	"tiktok-micro/app/api/internal/svc"
	"tiktok-micro/app/api/internal/types"
	"tiktok-micro/app/services/user/user_pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(req *types.LoginRequest) (resp *types.LoginResponse, err error) {
	res, err := l.svcCtx.UserRpc.Login(
		l.ctx, &user_pb.LoginRequest{
			Username: req.Username,
			Password: req.Password,
		})
	if err != nil {
		return nil, err
	}

	return &types.LoginResponse{
		StatusCode: res.StatusCode,
		StatusMsg:  res.StatusMsg,
		UserId:     res.UserId,
		Token:      res.Token,
	}, nil
}
