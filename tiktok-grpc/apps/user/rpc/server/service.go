package server

import (
	"context"
	"strconv"
	"tiktok-grpc/apps/user/pb"
	"tiktok-grpc/cmd"
	"tiktok-grpc/common/jwtx"
	"tiktok-grpc/common/redisx"
	"tiktok-grpc/common/utils"
	"tiktok-grpc/model"

	"golang.org/x/crypto/bcrypt"
)

type UserServiceServer struct {
	pb.UnimplementedUserServiceServer
}

func (p *UserServiceServer) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	db := cmd.DB

	if len(req.Username) > 30 || len(req.Password) > 30 {

		return &pb.RegisterResponse{

			StatusCode: 1,
			StatusMsg:  "长度超出限制!",
		}, nil
	}

	hashPassword, err := utils.HashPassword(req.Password) //hash加密密码
	if err != nil {
		return &pb.RegisterResponse{

			StatusCode: 1,
			StatusMsg:  "密码加密错误!",
		}, nil
	}

	var newUser model.User
	newUser.Name = req.Username
	newUser.Password = hashPassword
	result := db.Create(&newUser) //数据库创建用户

	//创建失败
	if result.Error != nil {
		return &pb.RegisterResponse{

			StatusCode: 1,
			StatusMsg:  "用户已存在!",
		}, nil
	}

	//创建成功
	token := jwtx.CreateUserToken(newUser.ID, newUser.Name) //获取token

	return &pb.RegisterResponse{

		StatusCode: 0,
		StatusMsg:  "注册成功!",

		UserId: newUser.ID,
		Token:  token,
	}, nil
}

func (p *UserServiceServer) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	db := cmd.DB
	var user model.User
	if len(req.Username) > 30 || len(req.Password) > 30 {

		return &pb.LoginResponse{

			StatusCode: 1,
			StatusMsg:  "长度超出限制!",
		}, nil
	}

	if db.Where("name = ?", req.Username).First(&user).RowsAffected == 0 {
		return &pb.LoginResponse{

			StatusCode: 1,
			StatusMsg:  "用户不存在!",
		}, nil
	}
	err := utils.ComparePassword(user.Password, req.Password)
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			// 密码不匹配
			return &pb.LoginResponse{

				StatusCode: 1,
				StatusMsg:  "用户密码错误!",
			}, nil
		} else {
			// 处理其他可能的错误
			return nil, err
		}
	}
	token := jwtx.CreateUserToken(user.ID, user.Name)

	return &pb.LoginResponse{
		StatusCode: 0,
		StatusMsg:  "登录成功!",
		UserId:     user.ID,
		Token:      token,
	}, nil
}
func (p *UserServiceServer) GetUser(ctx context.Context, req *pb.UserInfoRequest) (*pb.UserInfoResponse, error) {
	userID := req.UserId
	token := req.Token //获取前端注册用户ID和token
	userIDi, _ := strconv.Atoi(userID)
	//判断token是否有效
	tokenmsg, ok := jwtx.ParseToken(token)

	if int(tokenmsg.UserID) == userIDi && ok {

		user, err := redisx.UserCache(userID)
		if err != nil {

			return nil, err
		}

		return &pb.UserInfoResponse{

			StatusCode: 0,
			StatusMsg:  "登录成功!",

			User: utils.ConvertUserToProto(user),
		}, nil
	} else {
		//token无效
		return &pb.UserInfoResponse{

			StatusCode: 1,
			StatusMsg:  "token无效!",
		}, nil
	}

}
