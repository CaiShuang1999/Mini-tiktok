package service

import (
	"fmt"
	"tiktok-rest/middleware"
	"tiktok-rest/model"
	"tiktok-rest/utils"

	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type Response struct {
	// 状态码，0-成功，其他值-失败
	StatusCode int64 `json:"status_code"`
	// 返回状态描述
	StatusMsg string `json:"status_msg"`
}

type CreateUserReturn struct {
	Response
	// 用户鉴权token
	Token string `json:"token"`
	// 用户id
	UserID int64 `json:"user_id"`
}

type LoginUserReturn struct {
	Response
	UserId int64  `json:"user_id"`
	Token  string `json:"token"`
}

type UserInfoResponse struct {
	Response
	User *model.User `json:"user"` // 用户信息
}

func CreateUser(c *gin.Context) {

	userName := c.Query("username")
	userPassword := c.Query("password") //获取前端注册用户名和密码

	if len(userName) > 30 || len(userPassword) > 30 {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "用户名或密码长度超出限制"})
		return
	}

	db := utils.DB

	hashPassword, err := utils.HashPassword(userPassword) //hash加密密码
	if err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "密码不合法"})
		return
	}

	var newUser model.User
	newUser.Name = userName
	newUser.Password = hashPassword
	result := db.Create(&newUser) //数据库创建用户

	//创建失败
	if result.Error != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "用户已存在"})
		return
	}

	//创建成功
	token := middleware.CreateUserToken(newUser.ID, newUser.Name) //获取token
	backStatus := CreateUserReturn{
		Response: Response{
			StatusCode: 0,
			StatusMsg:  "success!",
		},
		UserID: newUser.ID,
		Token:  token,
	}
	c.JSON(http.StatusOK, backStatus) //返回响应

}

func LoginUser(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password") //获取前端注册用户名和密码

	if len(username) > 30 || len(password) > 30 {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "用户名或密码长度超出限制"})
		return
	}

	var user model.User
	db := utils.DB //数据库查询用户

	if db.Where("name = ?", username).First(&user).RowsAffected == 0 {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  "用户名或密码错误!",
		})
		return
	}

	err := utils.ComparePassword(user.Password, password)
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			// 密码不匹配
			c.JSON(http.StatusOK, Response{
				StatusCode: 1,
				StatusMsg:  "用户名或密码错误!",
			})
			return
		} else {
			// 处理其他可能的错误
			c.JSON(http.StatusOK, Response{
				StatusCode: 1,
				StatusMsg:  err.Error(),
			})
			return
		}
	}

	token := middleware.CreateUserToken(user.ID, user.Name)

	backStatus := LoginUserReturn{
		Response: Response{
			StatusCode: 0,
			StatusMsg:  "登录成功!",
		},
		UserId: user.ID,
		Token:  token,
	}

	c.JSON(http.StatusOK, backStatus)

}

func UserInfo(c *gin.Context) {
	userID := c.Query("user_id")
	token := c.Query("token") //获取前端注册用户ID和token
	userIDi, _ := strconv.Atoi(userID)

	//判断token是否有效
	tokenmsg, ok := middleware.ParseToken(token)

	if int(tokenmsg.UserID) == userIDi && ok {

		user, err := middleware.UserCache(userID)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		c.JSON(http.StatusOK, UserInfoResponse{
			Response: Response{StatusCode: 0, StatusMsg: "欢迎用户界面!"},
			User:     &user,
		})
	} else {
		//token无效
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "token无效"})
		return
	}

}
