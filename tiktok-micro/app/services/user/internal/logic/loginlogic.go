package logic

import (
	"context"

	"tiktok-micro/app/services/user/internal/svc"
	"tiktok-micro/app/services/user/user_pb"
	"tiktok-micro/common/cryptx"
	"tiktok-micro/common/jwtx"
	"tiktok-micro/model"

	"github.com/zeromicro/go-zero/core/logx"
	"golang.org/x/crypto/bcrypt"
)

type LoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *LoginLogic) Login(in *user_pb.LoginRequest) (*user_pb.LoginResponse, error) {
	// todo: add your logic here and delete this line

	username := in.Username
	password := in.Password //获取前端注册用户名和密码

	if len(username) > 30 || len(password) > 30 {
		return &user_pb.LoginResponse{
			StatusCode: 1, StatusMsg: "用户名或密码长度超出限制",
		}, nil
	}
	db := l.svcCtx.DB
	var user model.User

	if db.Where("name = ?", username).First(&user).RowsAffected == 0 {
		return &user_pb.LoginResponse{
			StatusCode: 1, StatusMsg: "用户名或密码错误!",
		}, nil
	}
	err := cryptx.ComparePassword(user.Password, password)
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			// 密码不匹配
			return &user_pb.LoginResponse{
				StatusCode: 1, StatusMsg: "用户名或密码错误!",
			}, nil
		} else {
			// 处理其他可能的错误
			return &user_pb.LoginResponse{
				StatusCode: 1, StatusMsg: err.Error(),
			}, nil
		}
	}

	token := jwtx.CreateUserToken(user.ID, user.Name)

	return &user_pb.LoginResponse{
		StatusCode: 0, StatusMsg: "登录成功!",
		UserId: user.ID, Token: token,
	}, nil
}
