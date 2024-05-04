package logic

import (
	"context"
	"strconv"
	"tiktok-micro/model"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var ctx = context.Background()

func CreateVideoInfoCache(db *gorm.DB, rdb *redis.Client, video_id string) (model.Video, error) {
	var video model.Video

	hashKey := "video:" + video_id

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

func VideoInfoCache(db *gorm.DB, rdb *redis.Client, video_id string) (model.Video, error) {
	var video model.Video

	hashKey := "video:" + video_id //视频key
	video_id_int, _ := strconv.Atoi(video_id)
	exists, err := rdb.Exists(ctx, hashKey).Result() //视频信息是否存在缓存
	if err != nil {
		return model.Video{}, err
	}

	if exists == 1 {
		//video缓存加载
		if err := rdb.HGetAll(ctx, hashKey).Scan(&video); err != nil {
			return model.Video{}, err
		}

	} else {
		video, err = CreateVideoInfoCache(db, rdb, video_id)
	}

	video.ID = int64(video_id_int)

	return video, err

}

func CreateVideoFavoritedCache(db *gorm.DB, rdb2 *redis.Client, video_id string) (int64, error) {
	var video model.Video

	var count int64 = 0

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

func VideoFavoritedCache(db *gorm.DB, rdb2 *redis.Client, video_id string) (int64, error) {

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
		count, err = CreateVideoFavoritedCache(db, rdb2, video_id)

	}
	return count, err

}

func VideoCache(db *gorm.DB, rdb *redis.Client, rdb2 *redis.Client, videoID int64) (model.Video, error) {

	video_id := strconv.Itoa(int(videoID))
	video, err := VideoInfoCache(db, rdb, video_id)
	if err != nil {
		return video, err
	}
	videoFavoriteCount, err := VideoFavoritedCache(db, rdb2, video_id)
	if err != nil {
		return video, err
	}

	video.FavoriteCount = videoFavoriteCount

	return video, err
}
