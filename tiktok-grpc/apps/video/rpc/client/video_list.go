package client

import (
	"log"
	"net/http"
	"tiktok-grpc/apps/video/video_pb"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func PubilshList(c *gin.Context) {
	serviceHost := viper.GetString("videoService.host")
	servicePort := viper.GetString("videoService.port")
	serviceAddress := serviceHost + ":" + servicePort
	conn, err := grpc.Dial(serviceAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := video_pb.NewVideoServiceClient(conn)
	req := &video_pb.PublishlistRequest{
		UserId: c.Query("user_id"),
		Token:  c.Query("token"),
	}
	res, err := client.PublishList(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status_code": 1, "status_msg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)

}
