package logic

import (
	"context"
	"math/rand"
	"strconv"
	"time"

	"tiktok-micro/app/services/user/internal/svc"
	"tiktok-micro/app/services/user/user_pb"
	"tiktok-micro/common/jwtx"
	"tiktok-micro/model"

	"github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type GetUserInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserInfoLogic {
	return &GetUserInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUserInfoLogic) GetUserInfo(in *user_pb.UserInfoRequest) (*user_pb.UserInfoResponse, error) {
	// todo: add your logic here and delete this line

	user_id := in.UserId
	token := in.Token
	//userIDi, _ := strconv.Atoi(user_id)
	//判断token是否有效
	_, ok := jwtx.ParseToken(token)

	if ok {
		user, err := UserCache(l.svcCtx.Config.Nginx.Addr,
			l.svcCtx.DB, l.svcCtx.Rdb_0, l.svcCtx.Rdb_1, user_id)

		if err != nil {
			return &user_pb.UserInfoResponse{StatusCode: 1, StatusMsg: err.Error()}, nil
		}

		return &user_pb.UserInfoResponse{StatusCode: 0, StatusMsg: "获取用户信息成功",
			User: &user_pb.UserInfo{

				Id:              user.ID,
				Name:            user.Name,
				FollowCount:     user.FollowCount,
				FollowerCount:   user.FollowerCount,
				Avatar:          user.Avatar,
				BackgroundImage: user.BackgroundImage,
				Signature:       user.Signature,
				TotalFavorited:  user.TotalFavorited,
				WorkCount:       user.WorkCount,
				FavoriteCount:   user.FavoriteCount,
			},
		}, nil
	} else {
		//token无效
		return &user_pb.UserInfoResponse{

			StatusCode: 1,
			StatusMsg:  "token无效!",
		}, nil
	}
}

var ctx = context.Background()

func CreateUserInfoCache(db *gorm.DB, rdb *redis.Client, user_id string) (model.User, error) {
	var user model.User

	hashKey := "user:" + user_id

	result := db.First(&user, user_id)
	if result.Error != nil {

		return model.User{}, result.Error

	}

	if err := rdb.HSet(ctx, hashKey, user).Err(); err != nil {
		return model.User{}, err
	}

	//设置过期随机时间
	expiration := time.Hour
	randomSeconds := rand.Intn(10 * 60)
	expirationWithRandomness := expiration + time.Duration(randomSeconds)*time.Second

	rdb.Expire(ctx, hashKey, expirationWithRandomness)

	return user, nil
}
func UserInfoCache(db *gorm.DB, rdb *redis.Client, user_id string) (model.User, error) {
	var user model.User

	hashKey := "user:" + user_id //用户信息key
	user_id_int, _ := strconv.Atoi(user_id)
	exists, err := rdb.Exists(ctx, hashKey).Result() //用户信息是否存在缓存
	if err != nil {
		return model.User{}, err
	}

	if exists == 1 {
		//用户缓存加载
		if err := rdb.HGetAll(ctx, hashKey).Scan(&user); err != nil {
			return model.User{}, err
		}

	} else {
		user, err = CreateUserInfoCache(db, rdb, user_id)
	}

	user.ID = int64(user_id_int)

	return user, err
}
func CreateUserFavoriteCache(db *gorm.DB, rdb2 *redis.Client, user_id string) (int64, error) {

	var user model.User

	var count int64 = 0

	result := db.First(&user, user_id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return count, result.Error
		} else {

			return count, result.Error
		}
	}

	userFavoriteCountKey := "favorite:user:" + user_id
	count = user.TotalFavorited

	rdb2.Set(ctx, userFavoriteCountKey, count, 0)

	return count, nil
}

func UserFavoriteCache(db *gorm.DB, rdb2 *redis.Client, user_id string) (int64, error) {

	var count int64 = 0
	userFavoriteCountKey := "favorite:user:" + user_id              //用户收到点赞数key
	exists2, err := rdb2.Exists(ctx, userFavoriteCountKey).Result() //用户收到点赞数key是否存在

	if err != nil {
		return count, err
	}

	if exists2 == 1 {
		//点赞数缓存加载

		userFavoriteCountStr, _ := rdb2.Get(ctx, userFavoriteCountKey).Result()
		userFavoriteCount, err := strconv.Atoi(userFavoriteCountStr)
		if err != nil {
			return count, err
		}
		count = int64(userFavoriteCount)
	} else {
		count, err = CreateUserFavoriteCache(db, rdb2, user_id)

	}
	return count, err
}
func UserCache(add string, db *gorm.DB, rdb *redis.Client, rdb2 *redis.Client, user_id string) (model.User, error) {
	//用户信息加载
	user, err := UserInfoCache(db, rdb, user_id)
	if err != nil {
		return user, err
	}
	//服务器资源定位
	user.Avatar = "http://" + add + user.Avatar
	user.BackgroundImage = "http://" + add + user.BackgroundImage

	//用户收到点赞数缓存加载
	userFavorited, err := UserFavoriteCache(db, rdb2, user_id)
	user.TotalFavorited = userFavorited
	return user, err

}
