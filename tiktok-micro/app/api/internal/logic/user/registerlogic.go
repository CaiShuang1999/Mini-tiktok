package user

import (
	"context"

	"tiktok-micro/app/api/internal/svc"
	"tiktok-micro/app/api/internal/types"
	"tiktok-micro/app/services/user/user_pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type RegisterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RegisterLogic) Register(req *types.RegisterRequest) (resp *types.RegisterResponse, err error) {
	// todo: add your logic here and delete this line
	res, err := l.svcCtx.UserRpc.Register(
		l.ctx, &user_pb.RegisterRequest{
			Username: req.Username,
			Password: req.Password,
		})
	if err != nil {
		return nil, err
	}

	return &types.RegisterResponse{
		StatusCode: res.StatusCode,
		StatusMsg:  res.StatusMsg,
		UserId:     res.UserId,
		Token:      res.Token,
	}, nil
}
