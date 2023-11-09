package logic

import (
	"context"
	"strconv"

	"tiktok-micro/app/services/favorite/favorite_pb"
	"tiktok-micro/app/services/favorite/internal/svc"
	"tiktok-micro/common/jwtx"
	"tiktok-micro/model"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type FavoriteActionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFavoriteActionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FavoriteActionLogic {
	return &FavoriteActionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FavoriteActionLogic) FavoriteAction(in *favorite_pb.FavoriteActionRequest) (*favorite_pb.FavoriteActionResponse, error) {
	// todo: add your logic here and delete this line
	ctx := context.Background()
	videoID := in.VideoId
	videoIDInt, _ := strconv.Atoi(videoID)
	token := in.Token
	action_type := in.ActionType

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

	db := l.svcCtx.DB

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

			if err := l.svcCtx.Rdb_0.Del(ctx, "user:"+useridstring).Err(); err != nil {
				return nil, err
			}
			if err := l.svcCtx.Rdb_0.Del(ctx, "user:"+authorIDstring).Err(); err != nil {
				return nil, err
			}
			if _, err := l.svcCtx.Rdb_1.Incr(ctx, "favorite:user:"+authorIDstring).Result(); err != nil {
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

			if err := l.svcCtx.Rdb_0.Del(ctx, "user:"+useridstring).Err(); err != nil {
				return nil, err
			}
			if err := l.svcCtx.Rdb_0.Del(ctx, "user:"+authorIDstring).Err(); err != nil {
				return nil, err
			}
			if _, err := l.svcCtx.Rdb_1.Decr(ctx, "favorite:user:"+authorIDstring).Result(); err != nil {
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
