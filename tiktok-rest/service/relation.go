package service

import (
	"net/http"
	"strconv"
	"tiktok-rest/middleware"
	"tiktok-rest/model"
	"tiktok-rest/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type FollowListResponse struct {
	Response
	Follows []model.User `json:"user_list"`
}

type FriendListResponse struct {
	// 状态码，0-成功，其他值-失败
	StatusCode string `json:"status_code"`
	// 返回状态描述
	StatusMsg string       `json:"status_msg"`
	Follows   []model.User `json:"user_list"`
}

func RelationAction(c *gin.Context) {
	token := c.Query("token")
	toUserIDStr := c.Query("to_user_id")
	toUserID, err := strconv.Atoi(toUserIDStr)
	if err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "ID无效"})
		return
	}
	actionType := c.Query("action_type")

	tokenMsg, ok := middleware.ParseToken(token)

	if !ok {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "token无效"})
		return
	}

	var relation model.Relation
	relation.UserID = tokenMsg.UserID
	relation.ToUserID = int64(toUserID)
	if relation.UserID == relation.ToUserID {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "不能对自己操作!"})
		return
	}

	db := utils.DB

	if db.Where("user_id=? AND to_user_id=?", relation.UserID, relation.ToUserID).Find(&model.Relation{}).RowsAffected == 0 {
		db.Create(&relation)
	}

	switch actionType {
	case "1":
		{

			db.Model(&model.Relation{}).Where("user_id=? AND to_user_id=?", relation.UserID, relation.ToUserID).
				Update("is_follow", true)
			db.Model(&model.User{}).Where("id = ?", relation.UserID).Update("follow_count", gorm.Expr("follow_count + ?", 1))
			db.Model(&model.User{}).Where("id = ?", relation.ToUserID).Update("follower_count", gorm.Expr("follower_count + ?", 1))
			c.JSON(http.StatusOK, Response{StatusCode: 0, StatusMsg: "关注成功!"})

		}

	case "2":
		{
			db.Model(&model.Relation{}).Where("user_id=? AND to_user_id=?", relation.UserID, relation.ToUserID).
				Update("is_follow", false)
			db.Model(&model.User{}).Where("id = ?", relation.UserID).Update("follow_count", gorm.Expr("follow_count - ?", 1))
			db.Model(&model.User{}).Where("id = ?", relation.ToUserID).Update("follower_count", gorm.Expr("follower_count - ?", 1))
			c.JSON(http.StatusOK, Response{StatusCode: 0, StatusMsg: "取消关注成功!"})
		}
	default:
		{
			c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "无效操作"})

		}

	}

}

func FollowList(c *gin.Context) {
	userID := c.Query("user_id")
	userIDInt, _ := strconv.Atoi(userID)

	token := c.Query("token")

	tokenMsg, ok := middleware.ParseToken(token)
	if int64(userIDInt) != tokenMsg.UserID && !ok {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "无效token"})
		return
	}

	var follows []model.User
	db := utils.DB

	err := db.Table("relation").
		Select("user.*").
		Joins("JOIN user ON user.id = relation.to_user_id").
		Where("relation.user_id = ? AND relation.is_follow=?", userID, true).
		Find(&follows).
		Error

	if err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "查找错误"})

		return
	}

	for i := range follows {
		follows[i].IsFollow = true
	}

	c.JSON(http.StatusOK, FollowListResponse{
		Response: Response{StatusCode: 0, StatusMsg: "查看关注列表"},
		Follows:  follows,
	})

}

func FansList(c *gin.Context) {
	userID := c.Query("user_id")
	userIDInt, _ := strconv.Atoi(userID)

	token := c.Query("token")

	tokenMsg, ok := middleware.ParseToken(token)
	if int64(userIDInt) != tokenMsg.UserID && !ok {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "无效token"})
		return
	}

	var fans []model.User
	db := utils.DB

	err := db.Table("relation").
		Select("user.*").
		Joins("JOIN user ON user.id = relation.user_id").
		Where("relation.to_user_id = ? AND is_follow=?", userID, true).
		Find(&fans).
		Error

	if err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "查找错误"})

		return
	}

	var count int64
	for i := range fans {

		db.Model(&model.Relation{}).Where("user_id = ? AND to_user_id = ? AND is_follow=?", userID, fans[i].ID, true).Count(&count)
		if count > 0 {
			fans[i].IsFollow = true
		}

	}

	c.JSON(http.StatusOK, FollowListResponse{
		Response: Response{StatusCode: 0, StatusMsg: "查看粉丝列表"},
		Follows:  fans,
	})

}

func FriendList(c *gin.Context) {
	userIDStr := c.Query("user_id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "无效ID"})
		return
	}

	token := c.Query("token")
	tokenMsg, ok := middleware.ParseToken(token)
	if tokenMsg.UserID != int64(userID) && !ok {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "无效token"})
		return
	}

	var friends []model.User

	db := utils.DB

	query := `
		SELECT user.*
		FROM relation 
		INNER JOIN user  ON relation .to_user_id = user.id
		WHERE relation.user_id = ? AND EXISTS (
			SELECT 1
			FROM relation
			WHERE user_id = user.id AND to_user_id = ?
		)
	`

	db.Raw(query, userID, userID).Scan(&friends)

	for i := range friends {
		friends[i].IsFollow = true
	}

	c.JSON(http.StatusOK, FriendListResponse{
		StatusCode: "0",
		StatusMsg:  "好友列表",
		Follows:    friends,
	})

}
