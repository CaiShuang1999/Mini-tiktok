package service

import (
	"fmt"
	"tiktok-rest/middleware"
	"tiktok-rest/utils"

	"tiktok-rest/model"

	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type MessageListResponse struct {
	// 状态码，0-成功，其他值-失败
	StatusCode string `json:"status_code"`
	// 返回状态描述
	StatusMsg   string          `json:"status_msg"`
	MessageList []model.Message `json:"message_list"`
}

func MessageAction(c *gin.Context) {

	token := c.Query("token")
	toUserIDstr := c.Query("to_user_id")
	actionType := c.Query("action_type")
	content := c.Query("content")

	if actionType != "1" {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "无效actionType"})
		return
	}

	tokenMsg, ok := middleware.ParseToken(token)
	if !ok {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "无效token"})
		return
	}

	toUserID, err := strconv.Atoi(toUserIDstr)
	if err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "无效ID"})
		return
	}

	var newMessage model.Message
	newMessage.FromUserID = tokenMsg.UserID
	newMessage.ToUserID = int64(toUserID)
	newMessage.Content = content
	newMessage.CreateTime = time.Now().Unix()

	db := utils.DB
	db.Create(&newMessage)

	c.JSON(http.StatusOK, Response{StatusCode: 0, StatusMsg: "发送信息成功"})
}

func MessageList(c *gin.Context) {

	token := c.Query("token")
	toUserIDstr := c.Query("to_user_id")
	preMsgTimeStr := c.Query("pre_msg_time")
	preMsgTime, _ := strconv.Atoi(preMsgTimeStr)

	tokenMsg, ok := middleware.ParseToken(token)
	if !ok {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "无效token"})
		return
	}

	toUserID, err := strconv.Atoi(toUserIDstr)
	if err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "无效ID"})
		return
	}

	var messageList []model.Message
	db := utils.DB
	db.Where("(from_user_id=? AND to_user_id=? AND create_time > ? ) OR (from_user_id=? AND to_user_id=? AND create_time > ? )",
		tokenMsg.UserID, toUserID, preMsgTime, toUserID, tokenMsg.UserID, preMsgTime).Find(&messageList)

	fmt.Println(messageList)
	c.JSON(http.StatusOK, MessageListResponse{
		StatusCode: "0", StatusMsg: "聊天列表",
		MessageList: messageList})
}
