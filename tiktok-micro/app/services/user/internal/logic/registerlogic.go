package logic

import (
	"context"

	"tiktok-micro/app/services/user/internal/svc"
	"tiktok-micro/app/services/user/user_pb"
	"tiktok-micro/common/cryptx"
	"tiktok-micro/common/jwtx"
	"tiktok-micro/model"

	"github.com/go-sql-driver/mysql"
	"github.com/zeromicro/go-zero/core/logx"
)

type RegisterLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RegisterLogic) Register(in *user_pb.RegisterRequest) (*user_pb.RegisterResponse, error) {
	// todo: add your logic here and delete this line

	userName := in.Username
	userPassword := in.Password //获取前端注册用户名和密码

	if len(userName) > 30 || len(userPassword) > 30 {

		return &user_pb.RegisterResponse{
			StatusCode: 1, StatusMsg: "用户名或密码长度超出限制",
		}, nil
	}

	db := l.svcCtx.DB

	hashPassword, err := cryptx.HashPassword(userPassword)
	if err != nil {
		return &user_pb.RegisterResponse{
			StatusCode: 1, StatusMsg: "密码不合法!",
		}, nil
	}

	var newUser model.User
	newUser.Name = userName
	newUser.Password = hashPassword
	result := db.Create(&newUser) //数据库创建用户

	//创建失败
	if result.Error != nil {
		if mysqlErr, ok := result.Error.(*mysql.MySQLError); ok && mysqlErr.Number == 1062 {
			return &user_pb.RegisterResponse{
				StatusCode: 1, StatusMsg: "用户名已存在!",
			}, nil
		} else {

			return &user_pb.RegisterResponse{
				StatusCode: 1, StatusMsg: err.Error(),
			}, nil
		}

	}

	//创建成功
	token := jwtx.CreateUserToken(newUser.ID, newUser.Name) //获取token

	return &user_pb.RegisterResponse{
		StatusCode: 0, StatusMsg: "创建用户成功!",
		UserId: newUser.ID, Token: token,
	}, nil
}
