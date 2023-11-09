package client

import (
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"tiktok-grpc/apps/video/video_pb"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type PublishQuery struct {
	Data  *multipart.FileHeader `form:"data"`
	Token string                `form:"token"`
	Title string                `form:"title"`
}

func UploadVideo(c *gin.Context) {
	token := c.PostForm("token")
	title := c.PostForm("title")

	fileHeader, err := c.FormFile("data")

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error:视频文件请求错误": err.Error()})
		return
	}
	file, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error:读取视频文件错误": err.Error()})
		return
	}
	defer file.Close()

	videoContent, err := io.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	serviceHost := viper.GetString("videoService.host")
	servicePort := viper.GetString("videoService.port")
	serviceAddress := serviceHost + ":" + servicePort
	conn, err := grpc.Dial(serviceAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := video_pb.NewVideoServiceClient(conn)
	req := &video_pb.UploadVideoRequest{
		Data:  []byte(videoContent),
		Title: title,
		Token: token,
	}
	res, err := client.UploadVideo(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status_code": 1, "status_msg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)

}
