package service

import (
	"net/http"
	"strconv"
	"tiktok-rest/middleware"
	"tiktok-rest/model"
	"tiktok-rest/utils"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

/*
	type CommentActionResponse struct {
		Response
		CommentText string `json:"comment"`
	}
*/
type CommentActionResponse struct {
	Response
	CommentText model.Comment `json:"comment"`
}
type CommentListResponse struct {
	Response
	Comments []model.Comment `json:"comment_list"`
}

func CommentAction(c *gin.Context) {
	token := c.Query("token")
	videoID := c.Query("video_id")
	actionType := c.Query("action_type")

	tokenmsg, ok := middleware.ParseToken(token)
	if !ok {
		//token无效
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "token无效"})
		return
	}

	videoIDInt, err := strconv.Atoi(videoID)
	if err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "无效video"})
		return
	}
	db := utils.DB

	var newComment model.Comment

	switch actionType {
	case "1":
		{

			newComment.CommentMsg = c.Query("comment_text")
			newComment.UserID = tokenmsg.UserID
			newComment.VideoID = int64(videoIDInt)
			newComment.CreateDate = time.Now().Format("06-01-02 15:04:05")

			db.Create(&newComment)

			db.Model(&model.Video{}).Where("id = ?", newComment.VideoID).Update("comment_count", gorm.Expr("comment_count  + ?", 1))

			c.JSON(http.StatusOK, CommentActionResponse{
				Response:    Response{StatusCode: 0, StatusMsg: "评论成功"},
				CommentText: newComment,
			})
		}
	case "2":
		{

			commentID := c.Query("comment_id")
			commentIDInt, _ := strconv.Atoi(commentID)

			db.First(&newComment, commentID)

			if newComment.UserID != tokenmsg.UserID {
				c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "无法删除别人的评论!"})
			}

			//软删除
			db.Model(&model.Comment{}).Where("id=?", int64(commentIDInt)).
				Update("delete_date", time.Now().Format("06-01-02 15:04:05"))

			db.Model(&model.Video{}).Where("id = ?", int64(videoIDInt)).Update("comment_count", gorm.Expr("comment_count  - ?", 1))
			c.JSON(http.StatusOK, CommentActionResponse{
				Response: Response{StatusCode: 0, StatusMsg: "删除评论成功"},
			})
		}
	default:
		{
			c.JSON(http.StatusOK, CommentActionResponse{
				Response: Response{StatusCode: 1, StatusMsg: "无效请求"},
			})
		}
	}

}

func CommentList(c *gin.Context) {
	token := c.Query("token")
	videoID := c.Query("video_id")

	_, ok := middleware.ParseToken(token)
	if !ok {
		//token无效
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "token无效"})
		return
	}

	db := utils.DB
	var comments []model.Comment

	db.Preload("UserMsg").Where("delete_date =? AND video_id=?", "", videoID).Order("id desc").Find(&comments)

	c.JSON(http.StatusOK, CommentListResponse{
		Response: Response{StatusCode: 0, StatusMsg: "评论区"},
		Comments: comments,
	})

}
