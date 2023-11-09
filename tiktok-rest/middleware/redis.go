package middleware

import (
	"context"
	"tiktok-rest/model"
	"tiktok-rest/utils"

	"math/rand"
	"strconv"
	"time"

	"gorm.io/gorm"
)

var ctx = context.Background()

func CreateUserInfoCache(user_id string) (model.User, error) {
	var user model.User
	var rdb = utils.RedisClient

	hashKey := "user:" + user_id

	db := utils.DB
	result := db.First(&user, user_id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return model.User{}, result.Error
		} else {

			return model.User{}, result.Error
		}
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

func UserInfoCache(user_id string) (model.User, error) {
	var user model.User
	var rdb = utils.RedisClient

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
		user, err = CreateUserInfoCache(user_id)
	}

	user.ID = int64(user_id_int)

	return user, err
}

func CreateUserFavoriteCache(user_id string) (int64, error) {

	var user model.User
	var rdb2 = utils.RedisClient2
	var count int64 = 0
	db := utils.DB
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

func UserFavoriteCache(user_id string) (int64, error) {
	var rdb2 = utils.RedisClient2
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
		count, err = CreateUserFavoriteCache(user_id)

	}
	return count, err
}

func UserCache(user_id string) (model.User, error) {
	//用户信息加载
	user, err := UserInfoCache(user_id)
	if err != nil {
		return user, err
	}
	//服务器资源定位
	user.Avatar = utils.StaticUrl + user.Avatar
	user.BackgroundImage = utils.StaticUrl + user.BackgroundImage

	//用户收到点赞数缓存加载
	userFavorited, err := UserFavoriteCache(user_id)
	user.TotalFavorited = userFavorited
	return user, err

}

func CreateVideoInfoCache(video_id string) (model.Video, error) {
	var video model.Video
	var rdb = utils.RedisClient

	hashKey := "video:" + video_id

	db := utils.DB
	result := db.First(&video, video_id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return model.Video{}, result.Error
		} else {

			return model.Video{}, result.Error
		}
	}

	if err := rdb.HSet(ctx, hashKey, video).Err(); err != nil {
		return model.Video{}, err
	}

	return video, nil
}

func VideoInfoCache(video_id string) (model.Video, error) {
	var video model.Video
	var rdb = utils.RedisClient

	hashKey := "video:" + video_id //视频key
	video_id_int, _ := strconv.Atoi(video_id)
	exists, err := rdb.Exists(ctx, hashKey).Result() //用户信息是否存在缓存
	if err != nil {
		return model.Video{}, err
	}

	if exists == 1 {
		//用户缓存加载
		if err := rdb.HGetAll(ctx, hashKey).Scan(&video); err != nil {
			return model.Video{}, err
		}

	} else {
		video, err = CreateVideoInfoCache(video_id)
	}

	video.ID = int64(video_id_int)

	return video, err

}

func CreateVideoFavoritedCache(video_id string) (int64, error) {
	var video model.Video
	var rdb2 = utils.RedisClient2
	var count int64 = 0
	db := utils.DB
	result := db.First(&video, video_id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return count, result.Error
		} else {

			return count, result.Error
		}
	}

	videoFavoriteCountKey := "favorite:video:" + video_id
	count = video.FavoriteCount

	rdb2.Set(ctx, videoFavoriteCountKey, count, 0)

	return count, nil
}

func VideoFavoritedCache(video_id string) (int64, error) {
	var rdb2 = utils.RedisClient2
	var count int64 = 0
	videoFavoriteCountKey := "favorite:video:" + video_id            //视频点赞数key
	exists2, err := rdb2.Exists(ctx, videoFavoriteCountKey).Result() //视频点赞数key是否存在

	if err != nil {
		return count, err
	}

	if exists2 == 1 {
		//点赞数缓存加载

		videoFavoriteCountStr, _ := rdb2.Get(ctx, videoFavoriteCountKey).Result()
		videoFavoriteCount, err := strconv.Atoi(videoFavoriteCountStr)
		if err != nil {
			return count, err
		}
		count = int64(videoFavoriteCount)
	} else {
		count, err = CreateVideoFavoritedCache(video_id)

	}
	return count, err

}

func VideoCache(videoID int64) (model.Video, error) {

	video_id := strconv.Itoa(int(videoID))
	video, err := VideoInfoCache(video_id)
	if err != nil {
		return video, err
	}
	videoFavoriteCount, err := VideoFavoritedCache(video_id)
	if err != nil {
		return video, err
	}
	video_author_id := strconv.Itoa(int(video.UserID))
	video.Author, err = UserCache(video_author_id)
	video.FavoriteCount = videoFavoriteCount

	return video, err
}
