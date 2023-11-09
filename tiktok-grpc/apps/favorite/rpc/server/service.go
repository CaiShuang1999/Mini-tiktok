package server

import (
	"context"
	"strconv"
	"tiktok-grpc/apps/favorite/favorite_pb"
	"tiktok-grpc/cmd"
	"tiktok-grpc/common/jwtx"
	"tiktok-grpc/common/utils"
	"tiktok-grpc/model"

	"gorm.io/gorm"
)

type favoriteServiceServer struct {
	favorite_pb.UnimplementedFavoriteServiceServer
}

func (p *favoriteServiceServer) FavoriteAction(ctx context.Context, req *favorite_pb.FavoriteActionRequest) (*favorite_pb.FavoriteActionResponse, error) {
	videoID := req.VideoId
	videoIDInt, _ := strconv.Atoi(videoID)
	token := req.Token
	action_type := req.ActionType

	//判断token是否有效
	tokenmsg, ok := jwtx.ParseToken(token)
	if !ok {
		//token无效

		return &favorite_pb.FavoriteActionResponse{
			StatusCode: 1, StatusMsg: "无效token",
		}, nil
	}

	var favorite model.Favorite
	userID := tokenmsg.UserID
	useridstring := strconv.Itoa(int(userID))

	db := cmd.DB

	favorite.UserID = userID
	favorite.VideoId = int64(videoIDInt)

	if db.Where("video_id=? AND user_id=?", favorite.VideoId, favorite.UserID).Find(&model.Favorite{}).RowsAffected == 0 {
		db.Create(&favorite)
	}

	var video model.Video
	db.Where("id=?", favorite.VideoId).Preload("Author").Find(&video)

	authorIDstring := strconv.Itoa(int(video.Author.ID))

	switch action_type {
	case "1":
		{

			db.Model(&model.Favorite{}).Where("video_id=? AND user_id=?", favorite.VideoId, favorite.UserID).
				Update("is_favorite", true)

			db.Model(&model.User{}).Where("id = ?", favorite.UserID).Update("favorite_count", gorm.Expr("favorite_count  + ?", 1))
			db.Model(&model.User{}).Where("id = ?", video.Author.ID).Update("total_favorited", gorm.Expr("total_favorited  + ?", 1))

			db.Model(&model.Video{}).Where("id = ?", favorite.VideoId).Update("favorite_count", gorm.Expr("favorite_count  + ?", 1))

			if err := cmd.RedisClient.Del(ctx, "user:"+useridstring).Err(); err != nil {
				return nil, err
			}
			if err := cmd.RedisClient.Del(ctx, "user:"+authorIDstring).Err(); err != nil {
				return nil, err
			}
			if _, err := cmd.RedisClient2.Incr(ctx, "favorite:user:"+authorIDstring).Result(); err != nil {
				return nil, err
			}

			return &favorite_pb.FavoriteActionResponse{
				StatusCode: 0, StatusMsg: "点赞成功",
			}, nil
		}

	case "2":
		{

			db.Model(&model.Favorite{}).Where("video_id=? AND user_id=?", favorite.VideoId, favorite.UserID).
				Update("is_favorite", false)

			db.Model(&model.User{}).Where("id = ?", favorite.UserID).Update("favorite_count", gorm.Expr("favorite_count  - ?", 1))
			db.Model(&model.User{}).Where("id = ?", video.Author.ID).Update("total_favorited", gorm.Expr("total_favorited  - ?", 1))
			db.Model(&model.Video{}).Where("id = ?", favorite.VideoId).Update("favorite_count", gorm.Expr("favorite_count  - ?", 1))

			if err := cmd.RedisClient.Del(ctx, "user:"+useridstring).Err(); err != nil {
				return nil, err
			}
			if err := cmd.RedisClient.Del(ctx, "user:"+authorIDstring).Err(); err != nil {
				return nil, err
			}
			if _, err := cmd.RedisClient2.Decr(ctx, "favorite:user:"+authorIDstring).Result(); err != nil {
				return nil, err
			}
			return &favorite_pb.FavoriteActionResponse{
				StatusCode: 0, StatusMsg: "取消点赞",
			}, nil
		}
	default:
		return &favorite_pb.FavoriteActionResponse{
			StatusCode: 1, StatusMsg: "无效操作",
		}, nil
	}

}
func (p *favoriteServiceServer) FavoriteList(ctx context.Context, req *favorite_pb.FavoriteListRequest) (*favorite_pb.FavoriteListResponse, error) {
	userID := req.UserId
	userIDint, _ := strconv.Atoi(userID)
	token := req.Token

	tokenmsg, ok := jwtx.ParseToken(token)
	if tokenmsg.UserID != int64(userIDint) && !ok {
		//token无效

		return &favorite_pb.FavoriteListResponse{
			StatusCode: 1, StatusMsg: "无效token",
		}, nil
	}

	db := cmd.DB
	var videos []model.Video

	err := db.Table("favorite").Select("video.*, user.*").
		Joins("JOIN video ON favorite.video_id = video.id").
		Joins("JOIN user ON video.user_id = user.id").
		Where("favorite.user_id = ? AND favorite.is_favorite=?", userID, true).
		Preload("Author"). // 预加载作者信息
		Find(&videos).Error

	if err != nil {

		return nil, err
	}
	for i := range videos {
		videos[i].IsFavorite = true
	}
	videos_msg := []*favorite_pb.VideoInfo{}
	for i := range videos {
		videos_msg = append(videos_msg, utils.ConvertVideoTofavoriteProto(videos[i]))
	}
	return &favorite_pb.FavoriteListResponse{
		StatusCode: 0, StatusMsg: "点赞列表",
		VideoList: videos_msg,
	}, nil
}
