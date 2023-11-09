package client

import (
	"log"
	"net/http"
	"tiktok-grpc/apps/favorite/favorite_pb"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func FavoriteAction(c *gin.Context) {
	serviceHost := viper.GetString("favoriteService.host")
	servicePort := viper.GetString("favoriteService.port")
	serviceAddress := serviceHost + ":" + servicePort
	conn, err := grpc.Dial(serviceAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	rpcclient := favorite_pb.NewFavoriteServiceClient(conn)
	req := &favorite_pb.FavoriteActionRequest{
		VideoId:    c.Query("video_id"),
		Token:      c.Query("token"),
		ActionType: c.Query("action_type"),
	}
	res, err := rpcclient.FavoriteAction(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}
