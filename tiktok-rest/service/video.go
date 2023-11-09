package service

import (
	"fmt"
	"tiktok-rest/middleware"
	"tiktok-rest/model"
	"tiktok-rest/utils"

	"log"
	"mime/multipart"
	"net/http"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type feedResponse struct {
	Response
	NextTime  int64         `json:"next_time"`  // 本次返回的视频中，发布最早的时间，作为下次请求时的latest_time
	VideoList []model.Video `json:"video_list"` // 视频列表
}

type PublishQuery struct {
	Data  *multipart.FileHeader `form:"data"`
	Token string                `form:"token"`
	Title string                `form:"title"`
}

type PublishList struct {
	Response
	VideoList *[]model.Video `json:"video_list"`
}

func FeedList(c *gin.Context) {

	token := c.Query("token")

	var videos []model.Video

	db := utils.DB
	db.Preload("Author").Order("id desc").Limit(30).Find(&videos)

	if token != "" {
		tokenmsg, ok := middleware.ParseToken(token)
		if !ok {
			//token无效
			c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "token无效"})
			return
		}

		userID := tokenmsg.UserID

		db := utils.DB

		for i := range videos {
			if db.Where("video_id = ? AND user_id = ? AND is_favorite=?", videos[i].ID, userID, true).Find(&model.Favorite{}).RowsAffected != 0 {
				videos[i].IsFavorite = true
			}

			if db.Where("user_id =? AND to_user_id=? AND is_follow=?", userID, videos[i].Author.ID, true).Find(&model.Relation{}).RowsAffected != 0 {
				videos[i].Author.IsFollow = true
			}

		}

	}
	for i := range videos {
		videos[i].Author.Avatar = utils.StaticUrl + videos[i].Author.Avatar
		videos[i].Author.BackgroundImage = utils.StaticUrl + videos[i].Author.BackgroundImage
	}
	c.JSON(http.StatusOK,
		feedResponse{
			Response:  Response{StatusCode: 0, StatusMsg: "视频流"},
			VideoList: videos,
		})

}

func PublishAPI(c *gin.Context) {
	var public PublishQuery
	if err := c.ShouldBind(&public); err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1,
			StatusMsg: "上传错误"})
		return
	}

	fileName := public.Data.Filename
	token := public.Token

	tokenmsg, ok := middleware.ParseToken(token)

	if !ok {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "token无效"})
		return
	}

	useridstring := strconv.Itoa(int(tokenmsg.UserID))

	timestamp := time.Now().UnixNano() / int64(time.Millisecond) // 获取当前时间的时间戳
	timestampString := strconv.Itoa(int(timestamp))
	fileNameS := strings.Split(fileName, ".")
	fileNameType := fileNameS[len(fileNameS)-1]
	savePath := "assets/videos/" + useridstring + "_" + timestampString + "." + fileNameType // 保存路径（用户ID+上传时间）

	err := c.SaveUploadedFile(public.Data, savePath)
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusInternalServerError, Response{StatusCode: 1, StatusMsg: "error:" + err.Error()})
		return
	}

	//封面
	inputFile := savePath
	outputFile := "assets/video_covers/" + useridstring + "_" + timestampString + ".jpg"
	timeOffset := "00:00:01"

	cmd := exec.Command("ffmpeg", "-i", inputFile, "-ss", timeOffset, "-vframes", "1", outputFile)

	err = cmd.Run()

	if err != nil {
		log.Fatal(err)
	}

	//上传成功后，video表数据库
	db := utils.DB
	var newVideo model.Video
	newVideo.UserID = tokenmsg.UserID
	newVideo.PlayURL = utils.StaticUrl + "/videos/" + useridstring + "_" + timestampString + "." + fileNameType
	newVideo.Title = public.Title
	newVideo.CoverURL = utils.StaticUrl + "/video_covers/" + useridstring + "_" + timestampString + ".jpg"
	newVideo.CreateTime = timestamp
	result := db.Create(&newVideo) //数据库创建video记录
	if result.Error != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "视频上传数据库错误!"})
		return
	}

	db.Model(&model.User{}).Where("id = ?", tokenmsg.UserID).Update("work_count", gorm.Expr("work_count + ?", 1))

	video, err := middleware.VideoCache(newVideo.ID)
	if err != nil {
		fmt.Println(video, err)
		return
	}

	c.JSON(http.StatusOK, Response{StatusCode: 0, StatusMsg: "上传成功!"})
}

func PublishListAPI(c *gin.Context) {

	userID := c.Query("user_id")

	var videos []model.Video
	db := utils.DB

	db.Preload("Author").Where("user_id = ?", userID).Find(&videos) //使用db.Preload("Author")预加载user表中"Author"信息，否则为0值

	c.JSON(http.StatusOK,
		PublishList{
			Response{StatusCode: 0, StatusMsg: "查询用户发布列表"},
			&videos})

}
