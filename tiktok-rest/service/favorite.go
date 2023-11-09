package service

import (
	"fmt"
	"tiktok-rest/middleware"
	"tiktok-rest/utils"

	"tiktok-rest/model"

	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func FavoriteAction(c *gin.Context) {
	videoID := c.Query("video_id")
	videoIDInt, _ := strconv.Atoi(videoID)
	token := c.Query("token")
	action_type := c.Query("action_type")

	//判断token是否有效
	tokenmsg, ok := middleware.ParseToken(token)
	if !ok {
		//token无效
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "token无效"})
		return
	}

	var favorite model.Favorite
	userID := tokenmsg.UserID
	db := utils.DB

	favorite.UserID = userID
	favorite.VideoId = int64(videoIDInt)

	if db.Where("video_id=? AND user_id=?", favorite.VideoId, favorite.UserID).Find(&model.Favorite{}).RowsAffected == 0 {
		db.Create(&favorite)
	}

	var video model.Video
	db.Where("id=?", favorite.VideoId).Preload("Author").Find(&video)

	switch action_type {
	case "1":
		{

			db.Model(&model.Favorite{}).Where("video_id=? AND user_id=?", favorite.VideoId, favorite.UserID).
				Update("is_favorite", true)

			db.Model(&model.User{}).Where("id = ?", favorite.UserID).Update("favorite_count", gorm.Expr("favorite_count  + ?", 1))
			db.Model(&model.User{}).Where("id = ?", video.Author.ID).Update("total_favorited", gorm.Expr("total_favorited  + ?", 1))

			db.Model(&model.Video{}).Where("id = ?", favorite.VideoId).Update("favorite_count", gorm.Expr("favorite_count  + ?", 1))

			c.JSON(http.StatusOK, Response{StatusCode: 0, StatusMsg: "点赞成功"})
		}

	case "2":
		{

			db.Model(&model.Favorite{}).Where("video_id=? AND user_id=?", favorite.VideoId, favorite.UserID).
				Update("is_favorite", false)

			db.Model(&model.User{}).Where("id = ?", favorite.UserID).Update("favorite_count", gorm.Expr("favorite_count  - ?", 1))
			db.Model(&model.User{}).Where("id = ?", video.Author.ID).Update("total_favorited", gorm.Expr("total_favorited  - ?", 1))
			db.Model(&model.Video{}).Where("id = ?", favorite.VideoId).Update("favorite_count", gorm.Expr("favorite_count  - ?", 1))
			c.JSON(http.StatusOK, Response{StatusCode: 0, StatusMsg: "取消点赞"})
		}
	default:
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "错误"})
	}

}

func FavoriteList(c *gin.Context) {
	userID := c.Query("user_id")
	userIDint, _ := strconv.Atoi(userID)
	token := c.Query("token")

	tokenmsg, ok := middleware.ParseToken(token)
	if tokenmsg.UserID != int64(userIDint) && !ok {
		//token无效
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "token无效"})
		return
	}

	db := utils.DB
	var videos []model.Video

	err := db.Table("favorite").Select("video.*, user.*").
		Joins("JOIN video ON favorite.video_id = video.id").
		Joins("JOIN user ON video.user_id = user.id").
		Where("favorite.user_id = ? AND favorite.is_favorite=?", userID, true).
		Preload("Author"). // 预加载作者信息
		Find(&videos).Error

	if err != nil {
		fmt.Println(err)
		return
	}
	for i := range videos {
		videos[i].IsFavorite = true
	}

	c.JSON(http.StatusOK, PublishList{Response: Response{StatusCode: 0, StatusMsg: "FAVORITE"},
		VideoList: &videos,
	})

}
